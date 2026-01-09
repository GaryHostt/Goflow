# Worker Pool Architecture

This document visualizes the bounded concurrency implementation in the iPaaS.

## The Problem: Unbounded Goroutines (POC)

```
┌────────────────────────────────────────────────────┐
│         Incoming Webhook Storm                      │
│  (1000 simultaneous requests)                       │
└────────────┬───────────────────────────────────────┘
             │
             ▼
   ┌─────────────────────────────────────────┐
   │  Goroutine Explosion (1000 spawned)     │
   │  ┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐┌──┐  │
   │  │G1││G2││G3││G4││G5││G6││G7││G8││G9│...│
   │  └──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘└──┘  │
   │  ... (990 more goroutines)              │
   └──────────────────┬──────────────────────┘
                      │
                      ▼
            ┌─────────────────────┐
            │  SQLite Database    │ ⚠️ LOCKED!
            │  (Single Writer)    │ (Database locked errors)
            └─────────────────────┘

Result: Server crash, database errors, OOM
```

---

## The Solution: Worker Pool (Production)

```
┌────────────────────────────────────────────────────┐
│         Incoming Webhook Storm                      │
│  (1000 simultaneous requests)                       │
└────────────┬───────────────────────────────────────┘
             │
             │ All requests queued
             ▼
   ┌─────────────────────────────────────────┐
   │        Job Queue (Buffered Channel)     │
   │  Capacity: 100 jobs                     │
   │  ┌────┬────┬────┬────┬────┬────┬────┐  │
   │  │Job1│Job2│Job3│Job4│... │Job99│100│  │
   │  └────┴────┴────┴────┴────┴────┴────┘  │
   │                                          │
   │  ⚠️ Queue Full? → Drop with warning     │
   └──────────────────┬──────────────────────┘
                      │
                      │ Pulled by workers
                      ▼
   ┌─────────────────────────────────────────┐
   │      Worker Pool (10 Fixed Workers)     │
   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐         │
   │  │W1 │ │W2 │ │W3 │ │W4 │ │W5 │         │
   │  └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘         │
   │    │     │     │     │     │            │
   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐         │
   │  │W6 │ │W7 │ │W8 │ │W9 │ │W10│         │
   │  └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘ └─┬─┘         │
   └────┼─────┼─────┼─────┼─────┼───────────┘
        │     │     │     │     │
        └─────┴─────┴─────┴─────┴─────────┐
                                           │
                      Max 10 concurrent    │
                      executions           │
                                           ▼
                            ┌─────────────────────┐
                            │  SQLite Database    │ ✅ SAFE!
                            │  (Max 10 writers)   │ (No lock errors)
                            └─────────────────────┘

Result: Predictable load, no crashes, graceful degradation
```

---

## Flow Diagram

```
┌──────────────┐
│ HTTP Request │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│ Submit to Pool   │───────┐
└──────────────────┘       │
                           │
                           ▼
                ┌──────────────────────┐
                │ Queue Full?          │
                └──────┬───────┬───────┘
                       │       │
                  YES  │       │ NO
                       │       │
                       ▼       ▼
              ┌─────────┐  ┌──────────┐
              │ Drop +  │  │ Add to   │
              │ Log Warn│  │ Queue    │
              └─────────┘  └────┬─────┘
                                │
                                ▼
                       ┌────────────────┐
                       │ Worker Pulls   │
                       │ Job from Queue │
                       └────┬───────────┘
                            │
                            ▼
                   ┌─────────────────────┐
                   │ Execute with Timeout│
                   │ (5 min max)         │
                   └────┬───────┬────────┘
                        │       │
                SUCCESS │       │ TIMEOUT/ERROR
                        │       │
                        ▼       ▼
                   ┌────────┐ ┌─────────┐
                   │ Log to │ │ Log     │
                   │ SQLite │ │ Failure │
                   └────────┘ └─────────┘
```

---

## Configuration

```go
// internal/engine/worker_pool.go

const (
    WorkerCount = 10          // Fixed number of worker goroutines
    QueueSize   = 100         // Buffered channel capacity
    JobTimeout  = 5 * Minute  // Max execution time per job
    QueueTimeout = 5 * Second // Wait time before dropping job
)
```

---

## Graceful Shutdown Sequence

```
                    SIGTERM Received
                           │
                           ▼
              ┌────────────────────────┐
              │ Close Job Queue        │ (No new jobs accepted)
              └────────────┬───────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │ Signal Workers to Stop │ (ctx.Cancel())
              └────────────┬───────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │ Wait for Workers       │ (sync.WaitGroup.Wait())
              │ - Max 30 seconds       │
              └────────────┬───────────┘
                           │
                  ┌────────┴────────┐
                  │                 │
            All Done         Timeout
                  │                 │
                  ▼                 ▼
         ┌────────────┐    ┌───────────────┐
         │ Shutdown   │    │ Force Shutdown│
         │ Complete ✅│    │ (In-flight    │
         └────────────┘    │ jobs dropped) │
                           └───────────────┘
```

---

## Performance Characteristics

| Scenario | POC Behavior | Worker Pool Behavior |
|----------|--------------|---------------------|
| **10 requests/sec** | 10 goroutines spawned | 10 workers idle, instant pickup |
| **100 requests/sec** | 100 goroutines spawned | 10 workers busy, 90 queued |
| **1000 requests/sec** | 1000 goroutines → crash | 10 workers busy, 100 queued, 890 dropped with warning |
| **Server shutdown** | In-flight jobs lost | 30-second drain, graceful completion |

---

## Monitoring Metrics

```go
// Available metrics from WorkerPool

workerPool.QueueLength()   // Current jobs waiting
workerPool.QueueCapacity() // Max queue size (100)

// Log output examples:
{
  "level": "info",
  "message": "Worker processing job",
  "worker_id": 3,
  "workflow_id": "wf_123",
  "queue_length": 45
}

{
  "level": "warn",
  "message": "Worker queue full, job dropped",
  "workflow_id": "wf_456",
  "queue_length": 100,
  "queue_cap": 100
}
```

---

## Scaling Decisions

### When to Increase Worker Count?

**Indicators:**
- "Queue full" warnings in logs
- Average queue length > 80% of capacity
- Execution latency increasing

**Action:**
```go
// Increase from 10 to 20 workers
pool := NewWorkerPool(20, logger)
```

### When to Increase Queue Size?

**Indicators:**
- Burst traffic patterns (e.g., 9am spike)
- Dropped jobs during known peak hours

**Action:**
```go
// Increase buffer from 100 to 500
jobQueue: make(chan WorkflowJob, 500)
```

### When to Move to Distributed Workers?

**Indicators:**
- Need > 50 workers
- Multiple API servers
- Geographic distribution
- Queue length consistently > 1000

**Action:**
- Replace in-process queue with Redis
- Deploy separate worker processes
- Use horizontal pod autoscaling (Kubernetes)

---

## Context Cancellation Integration

```
HTTP Request → Context Created
     │
     ▼
Submitted to Worker Pool (ctx passed)
     │
     ▼
Worker Pulls Job
     │
     ├─→ Check ctx.Done() ───┐
     │                        │ Cancelled?
     ▼                        ▼
Execute Workflow         Stop Immediately
     │
     ├─→ Check ctx.Done() ───┐
     │                        │ Cancelled?
     ▼                        ▼
Call Slack API          Stop Immediately
     │
     ├─→ Check ctx.Done() ───┐
     │                        │ Cancelled?
     ▼                        ▼
Save Log                 Don't Save (cancelled)
```

---

## SQLite Write Safety

**Without Worker Pool:**
```
┌──────────┐  ┌──────────┐  ┌──────────┐
│Goroutine1│  │Goroutine2│  │Goroutine3│
└────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │             │
     └─────────────┼─────────────┘
                   │
           All try to write simultaneously
                   ▼
            ┌─────────────┐
            │   SQLite    │
            │  (1 writer  │ ⚠️ "database is locked"
            │   at a time)│
            └─────────────┘
```

**With Worker Pool:**
```
┌──────────┐  ┌──────────┐  ┌──────────┐
│ Worker 1 │  │ Worker 2 │  │ Worker 3 │
└────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │             │
     │  (Sequential execution)   │
     │             │             │
     ▼             ▼             ▼
┌─────────────────────────────────┐
│        SQLite Database          │
│  Max 10 concurrent writes       │ ✅ Safe!
│  (Within SQLite's limits)       │
└─────────────────────────────────┘
```

---

## Production Checklist

- [x] **Worker Count Configured** (10 workers)
- [x] **Queue Size Set** (100 jobs buffer)
- [x] **Job Timeout Enforced** (5 minutes max)
- [x] **Queue Full Handling** (Log warning + drop)
- [x] **Graceful Shutdown** (30-second drain)
- [x] **Context Cancellation** (Stop on disconnect)
- [x] **Logging** (Worker ID, queue length, duration)
- [x] **Metrics Available** (QueueLength(), QueueCapacity())

---

**See Also:**
- [PRODUCTION_QUALITY.md](PRODUCTION_QUALITY.md) - Full architecture analysis
- [WHATS_NEW.md](WHATS_NEW.md) - v0.2.0 release notes
- [internal/engine/worker_pool.go](internal/engine/worker_pool.go) - Implementation

**Author**: Simple iPaaS Team  
**Date**: January 2026  
**Status**: Production-Ready ✅

