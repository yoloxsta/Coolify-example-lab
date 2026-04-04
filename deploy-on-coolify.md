# Deploy on Coolify — Step by Step

## Prerequisites

- Coolify is installed and running at `http://YourIP:8000`
- Your repo is pushed to GitHub (public or private)
- If private repo, connect GitHub to Coolify first

---

## Step 1 — Connect GitHub (skip if already done)

1. Go to Coolify dashboard → **Settings** → **Sources**
2. Click **Add Source** → choose **GitHub**
3. Follow the OAuth flow to authorize Coolify to access your repos

---

## Step 2 — Create a New Project

1. Go to **Projects** on the left sidebar
2. Click **+ Add** to create a new project
3. Give it a name like `myapp`
4. Click into the project, then click **+ New Resource**

---

## Step 3 — Add Docker Compose Resource

1. Choose **Docker Compose** as the resource type
2. Select your **server** (the one Coolify is running on)
3. Select **GitHub** as the source
4. Pick your repo (e.g. `your-username/myapp`)
5. Branch: `main`
6. Base Directory: `/` (if `docker-compose.yml` is at root of the repo)
7. Click **Continue**

---

## Step 4 — Configure Domains

In the resource's **General** settings page, you'll see two domain fields:

- **Domains for backend**: enter `http://api-myapp.YourIP.sslip.io`
- **Domains for frontend**: enter `http://myapp.YourIP.sslip.io`

Click **Save**.

---

## Step 5 — Update Frontend API URL

Before deploying, update `frontend/index.html` so the frontend knows where the backend is.

Replace this line:
```js
const API = window.location.origin.replace(/:\d+$/, '') + ':8080';
```

With:
```js
const API = 'http://api-myapp.YourIP.sslip.io';
```

Push the change:
```bash
git add .
git commit -m "set backend API url for production"
git push
```

---

## Step 6 — Deploy

1. Click **Redeploy** in Coolify
2. Watch the build logs — once all containers are started, you're live

---

## Step 7 — Verify

- Open `http://myapp.YourIP.sslip.io` in your browser
- You should see the Items page
- Try adding and deleting items
- Check backend health: `curl http://api-myapp.YourIP.sslip.io/api/health`

---

## Tips

- **Auto Deploy**: Enable **Webhooks** in Coolify so it auto-deploys on every `git push`
- **HTTPS**: Coolify + Traefik handles Let's Encrypt SSL automatically if you use a real domain
- **Logs**: Check container logs from Coolify dashboard if something goes wrong
- **Environment Variables**: You can override `DB_HOST`, `DB_PASSWORD`, etc. from Coolify's UI under each service's environment settings
- **No ports needed**: Don't use `ports` in `docker-compose.yml` — use `expose` instead. Coolify's Traefik proxy routes traffic via domains, and `ports` can conflict with Coolify's own services
