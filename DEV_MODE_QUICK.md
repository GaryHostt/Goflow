# ⚡ Quick Start: Dev Mode

## 🎯 One Command to Rule Them All

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/run_frontend_locally.sh
```

Then:
1. Open http://localhost:3000
2. Click **"Skip Login - Dev Mode"** (orange button with ⚡)
3. **Done!** You're in the dashboard

---

## 🚀 Full Setup (First Time)

```bash
# 1. Start backend
cd /Users/alex.macdonald/simple-ipass
docker compose up -d backend postgres elasticsearch

# 2. Start frontend with dev mode
./scripts/run_frontend_locally.sh

# 3. Open browser
open http://localhost:3000

# 4. Click "Skip Login - Dev Mode"
```

---

## 🔄 Daily Workflow

```bash
# Start backend (if not running)
docker compose up -d backend

# Start frontend
cd /Users/alex.macdonald/simple-ipass
./scripts/run_frontend_locally.sh

# Click "Skip Login - Dev Mode" → Start building! 🎉
```

---

## 🎨 What You Get

- ⚡ **Instant login** - No registration needed
- 🔄 **Hot reload** - Changes reflect immediately
- 🐛 **Easy debugging** - Console logs in terminal
- 🏃 **Fast iteration** - Build workflows quickly

---

## 🔐 Dev User Credentials

```
Email: dev@goflow.local
Password: dev123
```

*(Auto-created on first use)*

---

## 🛑 To Stop

Press **Ctrl+C** in the terminal running the frontend.

---

## 📖 Full Guide

See `DEV_MODE_GUIDE.md` for:
- Complete documentation
- Troubleshooting
- Security details
- Implementation details

---

**Happy coding!** 🚀⚡

