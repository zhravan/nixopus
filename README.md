<div align="center">

<a href="https://nixopus.com"><img width="1800" height="520" alt="Nixopus" src="https://github.com/user-attachments/assets/e103a9df-7abf-4f78-b75a-221331231247" /></a>

<h5 align="center">
  Vibe Deploy for full-stack apps — from code to live in under 60 seconds.
</h5>

<p align="center">
  <a href="https://nixopus.com"><b>Website</b></a> •
  <a href="https://docs.nixopus.ai"><b>Documentation</b></a> •
  <a href="https://nixopus.com/blog"><b>Blog</b></a> •
  <a href="https://discord.gg/skdcq39Wpv"><b>Discord</b></a> •
  <a href="https://github.com/raghavyuva/nixopus/discussions/262"><b>Roadmap</b></a>
</p>

<p align="center">
  <a href="https://github.com/raghavyuva/nixopus/stargazers"><img src="https://img.shields.io/github/stars/raghavyuva/nixopus?style=flat-square" alt="GitHub stars" /></a>
  <a href="https://github.com/raghavyuva/nixopus/network/members"><img src="https://img.shields.io/github/forks/raghavyuva/nixopus?style=flat-square" alt="GitHub forks" /></a>
  <a href="https://github.com/raghavyuva/nixopus/issues"><img src="https://img.shields.io/github/issues/raghavyuva/nixopus?style=flat-square" alt="GitHub issues" /></a>
  <a href="https://github.com/raghavyuva/nixopus/blob/master/LICENSE.md"><img src="https://img.shields.io/badge/license-FSL--1.1--ALv2-blue?style=flat-square" alt="License" /></a>
  <br>
  <a href="https://trendshift.io/repositories/15336" target="_blank"><img src="https://trendshift.io/api/badge/repositories/15336" alt="Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
</p>

</div>

---

<details>
<summary><h2>Table of Contents</h2></summary>

- [About the name](#about-the-name)
- [Getting started](#getting-started)
  - [Nixopus Cloud](#nixopus-cloud)
  - [Self-host](#self-host)
- [Features](#features)
- [Demo](#demo)
- [Self-hosting install](#self-hosting-install)
- [Links](#links)
- [Contributing](#contributing)
- [License](#license)

</details>

---

<a id="about-the-name"></a>
## About the name

Nixopus blends “octopus” (many arms, lots going on at once) with the Linux penguin vibe. It is **not** part of [NixOS](https://nixos.org/) or the Nix ecosystem.

---

<a id="getting-started"></a>
## Getting started

Nixopus is a deployment platform built to cut ops overhead: connect a repo, ship a build, get HTTPS and routing without hand-rolling CI or server sizing.

<a id="nixopus-cloud"></a>
### Nixopus Cloud

1. Create an account at [dashboard.nixopus.com](https://dashboard.nixopus.com/) (no card required).
2. Wait for your org’s isolated environment and `*.nixopus.ai` subdomain with HTTPS.
3. Connect GitHub, pick a repo, choose a build pack (Dockerfile or Docker Compose), set branch and env vars, deploy.

Details: [Quickstart](https://docs.nixopus.ai/getting-started/quickstart) · [Deploy your first app](https://docs.nixopus.ai/guides/deploying-apps)

<a id="self-host"></a>
### Self-host

Run the same stack on hardware you control: [Self-hosting](https://docs.nixopus.ai/getting-started/self-hosting) · [Installation](https://docs.nixopus.ai/self-hosting/installation) · [Configuration](https://docs.nixopus.ai/self-hosting/configuration)

> Watch releases on GitHub if you want notifications when we ship updates.

[![Star](assets/star.png)](https://github.com/raghavyuva/nixopus)

---

<a id="features"></a>
## Features

**Why Nixopus** (from the product docs):

- Zero ops — builds, deploys, SSL, and routing handled for you.
- Fast deploys — pick a build pack and go live in about a minute.
- Framework-agnostic — Next.js, Remix, Astro, FastAPI, Go, Rails, anything that runs in a container.
- Auto HTTPS — certificates via Caddy; custom domains supported.
- Rollbacks — revert from the Deployments UI.
- Built-in terminal — server shell in the dashboard; container shells from the Resources view.
- AI in the editor — VS Code / Cursor extension analyzes the repo, can generate Dockerfiles, and deploys from the sidebar.
- Open source — self-host or use cloud; API for automation.

**Guides**

| Topic | Doc |
| --- | --- |
| GitHub App & auto-deploy on push | [GitHub integration](https://docs.nixopus.ai/guides/github-integration) |
| Dockerfile / Compose, domains, env vars | [Deploy your first app](https://docs.nixopus.ai/guides/deploying-apps) |
| Compose-focused setup | [Docker Compose](https://docs.nixopus.ai/guides/docker-compose) |
| Env vars | [Environment variables](https://docs.nixopus.ai/guides/environment-variables) |
| Dashboard terminal | [Terminal](https://docs.nixopus.ai/guides/terminal) |
| Containers & images | [Container management](https://docs.nixopus.ai/guides/containers) |
| Marketplace extensions | [Extensions](https://docs.nixopus.ai/guides/extensions) |
| Metrics | [Charts](https://docs.nixopus.ai/guides/charts) |
| Alerts | [Notifications](https://docs.nixopus.ai/guides/notifications) |
| Uptime checks | [Health checks](https://docs.nixopus.ai/guides/health-checks) |
| In-dashboard AI | [AI chat](https://docs.nixopus.ai/guides/ai-chat) |
| Servers | [Server management](https://docs.nixopus.ai/guides/servers) |
| VS Code / Cursor | [Editor extension](https://docs.nixopus.ai/extension/overview) |

**Cloud (managed)**

Plans, billing, credits, API keys, custom domains, teams — [Cloud docs](https://docs.nixopus.ai/cloud/plans-and-billing).

**Concepts**

Architecture, auth, orgs, deployments, domains — [Core concepts](https://docs.nixopus.ai/concepts/architecture).

---

<a id="demo"></a>
## Demo

- **Flow**: connect GitHub → configure build → deploy → open the live URL (HTTPS).
- **Ops from the UI**: apps, domains, env vars, deployments, terminal, extensions.

### Screenshots (placeholders)

<div align="center">

| Dashboard | Hosting |
| --- | --- |
| ![Dashboard](assets/nixopus_dashboard.jpeg) | ![Hosting](assets/nixopus_dashboard.jpeg) |
| Deployments and activity | Repo → deploy |

| Terminal | Files |
| --- | --- |
| ![Terminal](assets/nixopus_dashboard.jpeg) | ![File manager](assets/nixopus_dashboard.jpeg) |
| Shell in the browser | Upload / edit / organize |

</div>

### Video (placeholder)

Add a demo link when you have one.

### Live preview (placeholder)

Add a public demo URL if you host one.

---

<a id="self-hosting-install"></a>
## Self-hosting install

Summarized from [Installation](https://docs.nixopus.ai/self-hosting/installation).

**Requirements**

| | Minimum |
| --- | --- |
| OS | Ubuntu 22.04+ (or Linux with systemd) |
| RAM | 2 GB |
| Docker | 20.10+ |
| Compose | v2+ |
| DNS | Domain or subdomain → server IP |

**Steps**

```bash
git clone https://github.com/nixopus/nixopus.git
cd nixopus
cp .env.self-hosted.sample .env
# Set at least: DATABASE_URL, BETTER_AUTH_SECRET, BETTER_AUTH_URL, ADMIN_EMAIL
docker compose up -d
```

That starts API, frontend, Postgres, Redis, and Caddy. Open the URL you set in `BETTER_AUTH_URL` and sign in with `ADMIN_EMAIL`. Run behind HTTPS in production; Caddy provisions certs when the public URL matches your DNS.

First successful sign-up is admin; after that, sign-up is closed and you invite users. Without an email provider, OTP codes show in Docker logs.

More: [Updating](https://docs.nixopus.ai/self-hosting/updating) · [Troubleshooting](https://docs.nixopus.ai/self-hosting/troubleshooting)

---

<a id="links"></a>
## Links

- **Site**: [nixopus.com](https://nixopus.com)
- **Docs**: [docs.nixopus.ai](https://docs.nixopus.ai)
- **API**: [API reference](https://docs.nixopus.ai/api-reference/introduction)
- **Discord**: [discord.gg/skdcq39Wpv](https://discord.gg/skdcq39Wpv)
- **Blog**: [nixopus.com/blog](https://nixopus.com/blog)
- **Roadmap**: [Discussions](https://github.com/raghavyuva/nixopus/discussions/262)
- **Issues**: [GitHub Issues](https://github.com/raghavyuva/nixopus/issues)

---

<a id="contributing"></a>
## Contributing

Issues and PRs are welcome. Start from [GitHub Issues](https://github.com/raghavyuva/nixopus/issues).

<a href="https://github.com/raghavyuva/nixopus/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=raghavyuva/nixopus" alt="Contributors" />
</a>

---

<a id="license"></a>
## License

FSL-1.1-ALv2 — see [LICENSE.md](./LICENSE.md).

---

<div align="center">

**Nixopus community**

[Back to top](#table-of-contents)

</div>
