# OAuth Setup Guide for Real-Time Forum

This guide will help you enable GitHub and Google OAuth login for your forum.

## 📋 Prerequisites
- GitHub account
- Google account
- Your app running on `http://localhost:3000` (frontend) and `http://localhost:8080` (backend)

---

## 🔧 Step 1: Register GitHub OAuth App

### 1.1 Go to GitHub Developer Settings
Visit: https://github.com/settings/developers

### 1.2 Click "New OAuth App"
Or use this direct link: https://github.com/settings/applications/new

### 1.3 Fill in the Application Details:
```
Application name: Real-Time Forum (or your preferred name)
Homepage URL: http://localhost:3000
Application description: A real-time forum with chat (optional)
Authorization callback URL: http://localhost:8080/api/auth/github/callback
```

⚠️ **IMPORTANT**: The callback URL must be exactly: `http://localhost:8080/api/auth/github/callback`

### 1.4 Click "Register Application"

### 1.5 Copy Your Credentials
You'll see:
- **Client ID** - Copy this
- Click "Generate a new client secret"
- **Client Secret** - Copy this (you won't see it again!)

### 1.6 Update .env File
Open `.env` file and add your credentials:
```bash
GITHUB_CLIENT_ID=your_actual_client_id_here
GITHUB_CLIENT_SECRET=your_actual_client_secret_here
```

---

## 🔧 Step 2: Register Google OAuth App

### 2.1 Go to Google Cloud Console
Visit: https://console.cloud.google.com

### 2.2 Create a New Project
1. Click the project dropdown at the top
2. Click "New Project"
3. Name it: "Real-Time Forum" (or your preferred name)
4. Click "Create"
5. Wait for project to be created, then select it

### 2.3 Enable Google+ API
1. In the left sidebar, go to: **APIs & Services** > **Library**
2. Search for "Google+ API"
3. Click on it
4. Click "Enable"

### 2.4 Configure OAuth Consent Screen
1. Go to: **APIs & Services** > **OAuth consent screen**
2. Choose "External" user type
3. Click "Create"
4. Fill in the required fields:
   ```
   App name: Real-Time Forum
   User support email: your-email@example.com
   Developer contact: your-email@example.com
   ```
5. Click "Save and Continue"
6. Skip "Scopes" - Click "Save and Continue"
7. Add test users (your email) - Click "Save and Continue"
8. Review and click "Back to Dashboard"

### 2.5 Create OAuth 2.0 Credentials
1. Go to: **APIs & Services** > **Credentials**
2. Click "+ CREATE CREDENTIALS" at the top
3. Select "OAuth client ID"
4. Choose Application type: "Web application"
5. Fill in details:
   ```
   Name: Real-Time Forum Web Client
   Authorized JavaScript origins: http://localhost:3000
   Authorized redirect URIs: http://localhost:8080/api/auth/google/callback
   ```

⚠️ **IMPORTANT**: The redirect URI must be exactly: `http://localhost:8080/api/auth/google/callback`

6. Click "Create"

### 2.6 Copy Your Credentials
You'll see a popup with:
- **Client ID** - Copy this
- **Client Secret** - Copy this

### 2.7 Update .env File
Open `.env` file and add your credentials:
```bash
GOOGLE_CLIENT_ID=your_actual_client_id_here
GOOGLE_CLIENT_SECRET=your_actual_client_secret_here
```

---

## 🚀 Step 3: Restart Your Application

### 3.1 Kill Running Servers
```bash
lsof -ti:8080,3000 | xargs kill -9
```

### 3.2 Start Servers with New Config
```bash
make dev
```

### 3.3 Test OAuth Login
1. Open browser to `http://localhost:3000`
2. Go to login page
3. Click "Continue with GitHub" or "Continue with Google"
4. Authorize the app
5. You should be redirected back and logged in!

---

## ✅ Verification Checklist

- [ ] Created GitHub OAuth app
- [ ] Copied GitHub Client ID to `.env`
- [ ] Copied GitHub Client Secret to `.env`
- [ ] Created Google Cloud project
- [ ] Enabled Google+ API
- [ ] Configured OAuth consent screen
- [ ] Created Google OAuth credentials
- [ ] Copied Google Client ID to `.env`
- [ ] Copied Google Client Secret to `.env`
- [ ] Restarted server with `make dev`
- [ ] Tested GitHub login
- [ ] Tested Google login

---

## 🐛 Troubleshooting

### OAuth redirect error
**Problem**: "Redirect URI mismatch" error
**Solution**: Double-check that callback URLs are EXACTLY:
- GitHub: `http://localhost:8080/api/auth/github/callback`
- Google: `http://localhost:8080/api/auth/google/callback`

### "Invalid client" error
**Problem**: Client ID or Secret is wrong
**Solution**: Copy credentials again from developer portals

### Still showing buttons but not working
**Problem**: `.env` file not loaded
**Solution**:
1. Check `.env` file exists: `ls -la .env`
2. Check credentials are filled (not empty)
3. Restart server: `make dev`

---

## 📝 Your Current .env File Location
File path: `/home/vapostol/real-time-forum/.env`

Edit with:
```bash
nano .env
# or
code .env
```

---

## 🎯 Next Steps After Setup

Once OAuth is working:
1. Users can register/login with GitHub or Google
2. Their profile data is automatically filled
3. No need to remember another password!

**Note**: For production deployment, you'll need to:
1. Update callback URLs to your production domain
2. Update `FRONTEND_BASE_URL` and `BACKEND_BASE_URL` in `.env`
3. Re-register OAuth apps with production URLs
