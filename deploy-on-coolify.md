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

### Option A — Using sslip.io (no domain needed)

- **Domains for backend**: `http://api-myapp.YourIP.sslip.io`
- **Domains for frontend**: `http://myapp.YourIP.sslip.io`

### Option B — Using your own domain (e.g. example.com with Cloudflare)

- **Domains for backend**: `https://api.example.com`
- **Domains for frontend**: `https://myapp.example.com`

See **Cloudflare DNS Setup** section below for how to point your domain.

Click **Save**.

---

## Step 5 — Cloudflare DNS Setup

If you have a wildcard domain on Cloudflare (e.g. `example.com`), follow these steps:

### 5.1 — Add DNS Records

Go to Cloudflare dashboard → your domain → **DNS** → **Records** → **Add Record**:

| Type | Name     | Content   | Proxy  |
|------|----------|-----------|--------|
| A    | myapp    | YourIP    | Proxied ☁️ |
| A    | api      | YourIP    | Proxied ☁️ |

This points `myapp.example.com` and `api.example.com` to your server.

> If you already have a wildcard record (`*.example.com` → YourIP), you don't need to add individual records. Any subdomain will automatically resolve to your server.

### 5.2 — SSL/TLS Settings in Cloudflare

Go to **SSL/TLS** → **Overview**:

- Set encryption mode to **Full** (not Full Strict)
- This is because Coolify/Traefik will handle the certificate on the server side, and Cloudflare handles it on the client side

### 5.3 — Enable SSL in Coolify

Coolify + Traefik auto-generates Let's Encrypt certificates. Just make sure:

1. In Coolify **Settings** → your server → check that **Wildcard Domain** is set (optional)
2. Use `https://` in your domain fields (Step 4 Option B)
3. Traefik will automatically request and renew SSL certs

---

## Step 6 — Update Frontend API URL

Update `frontend/index.html` so the frontend knows where the backend is.

Replace this line:
```js
const API = window.location.origin.replace(/:\d+$/, '') + ':8080';
```

With your backend domain:
```js
// If using sslip.io
const API = 'http://api-myapp.YourIP.sslip.io';

// If using your own domain
const API = 'https://api.example.com';
```

Push the change:
```bash
git add .
git commit -m "set backend API url for production"
git push
```

---

## Step 7 — Deploy

1. Click **Redeploy** in Coolify
2. Watch the build logs — once all containers are started, you're live

---

## Step 8 — Verify

- Open `https://myapp.example.com` (or `http://myapp.YourIP.sslip.io`) in your browser
- You should see the Items page
- Try adding and deleting items
- Check backend health: `curl https://api.example.com/api/health`

---

## Tips

- **Auto Deploy**: Enable **Webhooks** in Coolify so it auto-deploys on every `git push`
- **HTTPS**: Cloudflare handles client-side SSL, Coolify/Traefik handles server-side via Let's Encrypt
- **Logs**: Check container logs from Coolify dashboard if something goes wrong
- **Environment Variables**: You can override `DB_HOST`, `DB_PASSWORD`, etc. from Coolify's UI under each service's environment settings
- **No ports needed**: Don't use `ports` in `docker-compose.yml` — use `expose` instead. Coolify's Traefik proxy routes traffic via domains, and `ports` can conflict with Coolify's own services
- **Cloudflare Proxy**: Keep the orange cloud (Proxied) enabled for DDoS protection and caching. If you have SSL issues, try toggling to DNS Only (grey cloud) to debug
