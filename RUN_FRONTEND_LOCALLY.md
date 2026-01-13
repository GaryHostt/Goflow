# ✅ SOLUTION: Run Frontend Locally (Bypasses Docker Issues)

## Problem
Docker frontend is having connection issues. Running locally is faster and more reliable for development.

---

## 🚀 Run This Command

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/run_frontend_locally.sh
```

**This will:**
1. Create `.env.local` with correct API URL
2. Install npm dependencies (if needed)
3. Start the frontend development server

**Leave this terminal running!**

---

## ✅ Test It

1. **Open**: http://localhost:3000
2. **Register** with:
   - Email: `test@example.com`
   - Password: `password123`
3. **Should work!** ✅

---

## 🎯 Why This Works

- **Docker issue**: Frontend container has caching/networking problems
- **Local solution**: Runs directly on your machine, no Docker involved
- **Same result**: Connects to backend at `http://localhost:8080/api`

---

## 📊 What's Running

Now you have:
- ✅ **Backend** (in Docker): http://localhost:8080
- ✅ **Frontend** (local): http://localhost:3000
- ✅ **Postgres** (in Docker): localhost:5432
- ✅ **Elasticsearch** (in Docker): http://localhost:9200
- ✅ **Kibana** (in Docker): http://localhost:5601
- ✅ **Kong** (in Docker): http://localhost:8000-8002

**Only the frontend is running locally - everything else stays in Docker!**

---

## 🛑 To Stop

Press **Ctrl+C** in the terminal running the frontend.

---

## 🔄 To Switch Back to Docker Later

Once you're done testing and want to use Docker again:

```bash
# Stop local frontend (Ctrl+C)

# Start Docker frontend
docker compose up -d frontend
```

---

## 💡 Advantages of Running Locally

✅ **Faster**: No Docker build time  
✅ **Hot reload**: Changes reflect instantly  
✅ **Easier debugging**: See errors immediately in terminal  
✅ **No caching issues**: Always uses latest code  

---

## 🐛 If You See npm Errors

```bash
cd /Users/alex.macdonald/simple-ipass/frontend

# Clean install
rm -rf node_modules package-lock.json
npm install

# Run again
npm run dev
```

---

## ✅ Success Indicator

You should see:
```
- ready started server on 0.0.0.0:3000, url: http://localhost:3000
- event compiled client and server successfully
```

Then registration should work! 🎉

---

## 📚 Next Steps After It Works

1. ✅ **Register** a user
2. ✅ **Login** with that user
3. ✅ **Create a workflow** (Dashboard → Workflows → New)
4. ✅ **Add credentials** (Dashboard → Connections)
5. ✅ **View logs** (Dashboard → Logs)

---

**Run the script now and let me know if it works!** 🚀

