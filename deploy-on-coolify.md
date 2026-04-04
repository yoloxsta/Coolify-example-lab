# Deploy on Coolify — Step by Step

## Prerequisites

- Coolify is installed and running (e.g. `http://YOUR_SERVER_IP:8000`)
- Your repo is pushed to GitHub (public or private)
- If private repo, you need to connect GitHub to Coolify first

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

## Step 4 — Configure Services

Coolify will detect 3 services from `docker-compose.yml`: **frontend**, **backend**, **db**.

### Frontend
- Click on the **frontend** service
- Under **Domains**, add your domain or use Coolify's generated URL
  - Example: `http://myapp.YOUR_SERVER_IP.sslip.io`
- Port: `80` (nginx inside the container)

### Backend
- Click on the **backend** service
- Under **Domains**, add a domain for the API
  - Example: `http://api.myapp.YOUR_SERVER_IP.sslip.io`
- Port: `8080`

### Database (db)
- No domain needed — it's internal only
- Coolify will handle the volume (`pgdata`) automatically

---

## Step 5 — Deploy

1. Click **Deploy** at the top
2. Coolify will pull the repo, build the Docker images, and start all containers
3. Watch the build logs — if everything is green, you're live

---

## Step 6 — Update Frontend API URL

After deploying with domains, the frontend JS auto-detection (`localhost:8080`) won't work with real domains. You need to update `frontend/index.html`:

Replace this line:
```js
const API = window.location.origin.replace(/:\d+$/, '') + ':8080';
```

With your actual backend domain:
```js
const API = 'http://api.myapp.YOUR_SERVER_IP.sslip.io';
```

Or for HTTPS:
```js
const API = 'https://api.myapp.yourdomain.com';
```

Then push the change and redeploy from Coolify.

---

## Step 7 — Verify

- Open your frontend domain in a browser
- You should see the Items page
- Try adding and deleting items
- Check backend health: `curl http://api.myapp.YOUR_SERVER_IP.sslip.io/api/health`

---

## Tips

- **Auto Deploy**: In Coolify resource settings, enable **Webhooks** so it auto-deploys on every `git push`
- **HTTPS**: Coolify + Traefik handles Let's Encrypt SSL automatically if you use a real domain
- **Logs**: Check container logs from Coolify dashboard if something goes wrong
- **Environment Variables**: You can override `DB_HOST`, `DB_PASSWORD`, etc. from Coolify's UI under each service's environment settings
