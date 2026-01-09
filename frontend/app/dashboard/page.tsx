'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { workflows, logs } from '@/lib/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

interface Workflow {
  id: string
  name: string
  is_active: boolean
}

interface Log {
  id: string
  workflow_name: string
  status: string
  message: string
  executed_at: string
}

export default function DashboardPage() {
  const [workflowsList, setWorkflowsList] = useState<Workflow[]>([])
  const [logsList, setLogsList] = useState<Log[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [workflowsData, logsData] = await Promise.all([
        workflows.list(),
        logs.list(),
      ])
      setWorkflowsList(workflowsData || [])
      setLogsList(logsData || [])
    } catch (err) {
      console.error('Failed to load dashboard data', err)
    } finally {
      setLoading(false)
    }
  }

  const activeWorkflows = workflowsList.filter(w => w.is_active).length
  const recentLogs = logsList.slice(0, 5)
  const successRate = logsList.length > 0
    ? Math.round((logsList.filter(l => l.status === 'success').length / logsList.length) * 100)
    : 0

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome to your iPaaS platform
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-6 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardDescription>Total Workflows</CardDescription>
            <CardTitle className="text-4xl">{workflowsList.length}</CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardDescription>Active Workflows</CardDescription>
            <CardTitle className="text-4xl">{activeWorkflows}</CardTitle>
          </CardHeader>
        </Card>

        <Card>
          <CardHeader>
            <CardDescription>Success Rate</CardDescription>
            <CardTitle className="text-4xl">{successRate}%</CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>Get started with your integrations</CardDescription>
        </CardHeader>
        <CardContent className="flex gap-4">
          <Link href="/dashboard/workflows/new">
            <Button>Create Workflow</Button>
          </Link>
          <Link href="/dashboard/connections">
            <Button variant="outline">Connect Service</Button>
          </Link>
        </CardContent>
      </Card>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>Latest workflow executions</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>Loading...</p>
          ) : recentLogs.length === 0 ? (
            <p className="text-muted-foreground text-center py-4">
              No activity yet. Create a workflow to get started!
            </p>
          ) : (
            <div className="space-y-4">
              {recentLogs.map((log) => (
                <div key={log.id} className="flex items-center justify-between border-b pb-3">
                  <div>
                    <p className="font-medium">{log.workflow_name}</p>
                    <p className="text-sm text-muted-foreground">{log.message}</p>
                  </div>
                  <Badge variant={log.status === 'success' ? 'success' : 'destructive'}>
                    {log.status}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

