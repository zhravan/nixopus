<div align="center">

<a href="https://nixopus.com"><img width="1800" height="520" alt="Nixopus" src="https://github.com/user-attachments/assets/e103a9df-7abf-4f78-b75a-221331231247" /></a>

<h5 align="center">
  Vibe Deploy for full-stack apps. From code to live in under 60 seconds.
</h5>

<p align="center">
  <a href="https://nixopus.com"><b>Website</b></a> •
  <a href="https://docs.nixopus.ai/"><b>Documentation</b></a> •
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

<a id="table-of-contents"></a>

<details>
<summary><h2>Table of Contents</h2></summary>

- [About the Name](#about-the-name)
- [Getting Started](#getting-started)
- [Why Nixopus](#why-nixopus)
- [Explore](#explore)
- [Demo](#demo)
- [Installation](#installation)
- [Links](#links)
- [Contributing](#contributing)
- [License](#license)

</details>

---

<a id="about-the-name"></a>
## About the Name

Nixopus blends “octopus” (many arms, lots going at once) with the Linux penguin vibe. Despite the name, it is not part of [NixOS](https://nixos.org/) or the Nix ecosystem.

---

<a id="getting-started"></a>
## Getting Started

Nixopus is a deployment platform built to strip out busywork: connect a repository, deploy, and get HTTPS without wiring your own CI or sizing servers.

**Nixopus Cloud** — Create an account, get an isolated environment with a `*.nixopus.ai` URL, link GitHub, pick a build pack (Dockerfile or Docker Compose), set branch and env vars, deploy. Typical turnaround is under a minute. Details: [Quickstart](https://docs.nixopus.ai/getting-started/quickstart).

**Self-hosted** — Same idea on hardware you control. See [Self-hosting](https://docs.nixopus.ai/getting-started/self-hosting) and [Installation](#installation) below.

> [!TIP]
> Star the repo if you want GitHub release notifications.

[![Star](assets/star.png)](https://github.com/raghavyuva/nixopus)

---

<a id="why-nixopus"></a>
## Why Nixopus

- **Zero ops** — Push code, get a live URL with HTTPS. No CI pipelines, no server sizing.
- **60-second deploys** — Connect your repo, pick a build pack, and go live.
- **Framework agnostic** — Next.js, Django, Rails, Go, or anything that runs in a container.
- **Built-in terminal** — SSH into your server or containers from the dashboard.
- **Open source** — Self-host on your own machine or use Nixopus Cloud.
- **AI-powered** — Deploy from VS Code or Cursor with an agent that helps with configs and shipping.

---

<a id="explore"></a>
## Explore

| Topic | Description |
| --- | --- |
| [Deploy your first app](https://docs.nixopus.ai/guides/deploying-apps) | Repo to production in under 60 seconds. |
| [GitHub integration](https://docs.nixopus.ai/guides/github-integration) | Connect repos and auto-deploy on push. |
| [Custom domains](https://docs.nixopus.ai/cloud/custom-domains) | Bring your domain; HTTPS handled for you. |

More guides (terminal, Compose, env vars, notifications, health checks, extensions, etc.) live in the [Guides](https://docs.nixopus.ai/) section of the docs.

---

<a id="demo"></a>
## Demo

Walkthrough of deploys and day-to-day use is still being captured here.

### What you’ll see

- Connect a GitHub repo and ship with Dockerfile or Compose.
- Operate from the dashboard: domains, env vars, deployments, logs.
- Use the in-product terminal and other dashboard tools.

### Screenshots (placeholders)

<div align="center">

| Dashboard | Hosting |
| --- | --- |
| ![Dashboard Overview](assets/nixopus_dashboard.jpeg) | ![Hosting](assets/nixopus_dashboard.jpeg) |
| *Deployments and activity* | *Repo connect and releases* |

| Terminal | Files |
| --- | --- |
| ![Terminal](assets/nixopus_dashboard.jpeg) | ![File manager](assets/nixopus_dashboard.jpeg) |
| *Shell from the browser* | *Upload and edit on the server* |

</div>

### Video (placeholder)

- Demo link or file: TBD.

### Live preview (placeholder)

- Hosted demo URL: TBD.

---

<a id="installation"></a>
## Installation

Self-host the full stack on a box you own. Requirements from the docs:

| Requirement | Minimum |
| --- | --- |
| OS | Ubuntu 22.04+ (or Linux with systemd) |
| RAM | 2 GB |
| Docker | 20.10+ |
| Docker Compose | v2+ |
| DNS | Domain or subdomain pointing at the server |

```bash
git clone https://github.com/nixopus/nixopus.git
cd nixopus
cp .env.self-hosted.sample .env
```

Set at least:

| Variable | Purpose |
| --- | --- |
| `DATABASE_URL` | PostgreSQL connection string |
| `BETTER_AUTH_SECRET` | Random 32+ character signing secret |
| `BETTER_AUTH_URL` | Public URL of your install |
| `ADMIN_EMAIL` | Initial admin email |

```bash
docker compose up -d
```

Then open the URL you set in `BETTER_AUTH_URL` and sign in with `ADMIN_EMAIL`. Run behind HTTPS in production; Caddy renews certs when a domain points at the server and `BETTER_AUTH_URL` matches.

Full variable list and behavior: [Configuration](https://docs.nixopus.ai/self-hosting/configuration). Operational notes: [Installation](https://docs.nixopus.ai/self-hosting/installation), [Updating](https://docs.nixopus.ai/self-hosting/updating), [Troubleshooting](https://docs.nixopus.ai/self-hosting/troubleshooting).

---

<a id="links"></a>
## Links

| Resource | URL |
| --- | --- |
| Website | [nixopus.com](https://nixopus.com) |
| Documentation | [docs.nixopus.ai](https://docs.nixopus.ai/) |
| API reference | [docs.nixopus.ai/api-reference](https://docs.nixopus.ai/api-reference/introduction) |
| Editor extension | [docs.nixopus.ai/extension](https://docs.nixopus.ai/extension/overview) |
| Discord | [discord.gg/skdcq39Wpv](https://discord.gg/skdcq39Wpv) |
| Blog | [nixopus.com/blog](https://nixopus.com/blog) |
| Roadmap | [Discussions](https://github.com/raghavyuva/nixopus/discussions/262) |
| Issues | [GitHub Issues](https://github.com/raghavyuva/nixopus/issues) |

---

<a id="contributing"></a>
## Contributing

Issues and pull requests are welcome. Start from [GitHub Issues](https://github.com/raghavyuva/nixopus/issues) if you are unsure where to jump in.

<a href="https://github.com/raghavyuva/nixopus/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=raghavyuva/nixopus" alt="Nixopus contributors" />
</a>

---

<a id="license"></a>
## License

FSL-1.1-ALv2. See [LICENSE.md](./LICENSE.md).

---

<div align="center">

**Made with ❤️ by the Nixopus community**

[Back to top](#table-of-contents)

</div>
