# ⚡ GoFlow Test Suite - Quick Reference

## 🚀 One-Line Test Execution

```bash
cd /Users/alex.macdonald/simple-ipass && ./scripts/run_all_tests.sh
```

---

## 📋 Individual Commands

```bash
# Start platform
docker compose up -d && sleep 60

# Configure Kong
./scripts/configure_kong_elk.sh

# Test connectors (30s)
go run scripts/connector_test.go

# Test Kong (45s)
go run scripts/kong_test.go

# View results
open http://localhost:5601
```

---

## 🎯 Makefile Shortcuts

```bash
make docker-up              # Start platform
make configure-kong-elk     # Setup Kong logs
make test-connectors        # Test 18 connectors
make test-kong              # Test 5 Kong patterns
make test-full              # Run everything
```

---

## 📊 Expected Results

| Test | Pass | Skip | Fail | Duration |
|------|------|------|------|----------|
| Connectors | 12 | 4 | 0 | 30s |
| Kong | 5 | 0 | 0 | 45s |
| **Total** | **17** | **4** | **0** | **90s** |

---

## 🔗 Quick Links

- **Kibana**: http://localhost:5601
- **Kong Admin**: http://localhost:8001
- **Kong Manager**: http://localhost:8002
- **Backend API**: http://localhost:8080
- **Frontend**: http://localhost:3000

---

## 🐛 Quick Troubleshooting

**Services not starting?**
```bash
docker compose down && docker compose up -d
```

**Tests failing?**
```bash
docker compose ps  # Check health
docker compose logs backend  # View logs
```

**Kibana not showing data?**
```bash
curl http://localhost:9200/_cat/indices?v  # Check indices
```

---

## ✅ Success Indicators

- ✅ "All critical tests passed!"
- ✅ "All Kong Gateway tests passed!"
- ✅ Test results visible in Kibana

---

**See `HOW_TO_RUN_TESTS.md` for detailed instructions!**

