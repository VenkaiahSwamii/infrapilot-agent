# LinkedIn Post Series — Building InfraPilot Enterprise

Six posts documenting the engineering journey. Post roughly 2-4 days apart. Each includes a suggested image/attachment and hashtags.

---

## Post 1 — Building the Monitoring Agent

**Suggested attachment:** Screenshot of agent terminal output enrolling and streaming metrics.

> 🚀 I've been building InfraPilot Enterprise — a production-style infrastructure monitoring platform — and today I want to talk about the piece that started it all: the monitoring agent.
>
> It's a cross-platform Go binary that runs on the machines you want to monitor. A few design goals I set for it:
>
> ✅ Under 1% CPU and ~20MB RAM overhead — it should never be the thing that causes an incident
> ✅ Cross-compile cleanly to Linux, Windows, and macOS from a single codebase
> ✅ Pluggable architecture — system metrics, Docker, Kubernetes, logs, and security checks are all separate plugins behind a common interface
> ✅ Secure enrollment — agents authenticate with a one-time token, then communicate via API key + WebSocket
>
> Go's standard library made cross-platform system metrics collection (CPU, memory, disk, network, processes) surprisingly manageable, and its concurrency model is a perfect fit for an agent juggling multiple collectors on independent tickers.
>
> Next post: how the agent talks to Docker and streams container metrics in real time. 🐳
>
> #golang #devops #infrastructure #softwareengineering #buildinpublic

---

## Post 2 — Implementing Docker Monitoring

**Suggested attachment:** Screenshot of the Docker tab showing live container metrics.

> 🐳 Part 2 of building InfraPilot Enterprise: Docker monitoring.
>
> Once the base agent was collecting system metrics, the natural next step was containers. I integrated directly with the Docker Engine API to:
>
> - Enumerate running containers and their resource usage (CPU, memory)
> - Track container lifecycle events (created, started, stopped, removed)
> - Expose start/stop/restart actions back through the dashboard
> - Stream container metrics in near real time alongside host-level metrics
>
> The interesting challenge here wasn't the API integration itself — it was normalizing Docker's stats format (which is delta-based and a little awkward) into the same time-series shape as host metrics, so the dashboard could treat them uniformly.
>
> This is the kind of detail that doesn't show up in a demo video but matters a lot for a platform that's supposed to feel unified rather than bolted together.
>
> Next: bringing Kubernetes into the picture. ☸️
>
> #docker #golang #devops #containers #buildinpublic

---

## Post 3 — Kubernetes Integration

**Suggested attachment:** Screenshot of the Kubernetes tab (pods/nodes) or a terminal running `kubectl get pods`.

> ☸️ Part 3: adding Kubernetes support to InfraPilot Enterprise.
>
> Docker monitoring covers standalone hosts, but a lot of real infrastructure today runs on Kubernetes. So I added a K8s plugin to the agent that talks to the cluster API (via client-go) to surface:
>
> - Pod status, restarts, and resource usage
> - Node health and capacity
> - Deployment status
>
> The biggest lesson here: Kubernetes' API is powerful but verbose. I spent real time deciding what to expose versus what would just be noise on a dashboard meant for quick operational awareness rather than being a full kubectl replacement. Less can genuinely be more when the audience is "is anything on fire right now?"
>
> This also pushed me to formalize the agent's plugin interface so that Linux, Docker, and Kubernetes collectors all implement the same contract — which made testing and extending it much easier.
>
> Next: the real-time architecture that gets all this data to the browser in milliseconds. ⚡
>
> #kubernetes #golang #devops #cloudnative #buildinpublic

---

## Post 4 — WebSocket Architecture

**Suggested attachment:** A simple architecture diagram (Agent → Backend → Event Bus → WebSocket Hub → Browser).

> ⚡ Part 4: how InfraPilot Enterprise gets data to the browser in real time.
>
> Polling is simple but doesn't scale well and adds latency. Instead, I built:
>
> - An internal **event bus** in the Go backend — metrics, alerts, and machine state changes are all published as events
> - Multiple **subscribers**: one persists events to PostgreSQL, one pushes to a WebSocket hub, one evaluates alert rules, one writes audit logs
> - A **WebSocket hub** that manages per-client subscriptions so each connected dashboard only receives the events relevant to what it's viewing
>
> This decoupled design means adding a new "reaction" to an event (say, a new notification channel) doesn't touch the metric ingestion path at all — it's just a new subscriber on the bus.
>
> The result: dashboard updates in well under 100ms from the moment an agent reports a metric, even under load. Backpressure handling and heartbeat-based reconnection were the trickiest parts to get right — WebSocket connections fail in more interesting ways than you'd expect.
>
> Next: taking this from "runs on my laptop" to a real production deployment. 🚀
>
> #softwarearchitecture #websocket #golang #systemdesign #buildinpublic

---

## Post 5 — Production Deployment

**Suggested attachment:** Screenshot of `kubectl get pods -n infrapilot` or the live dashboard URL.

> 🚀 Part 5: taking InfraPilot Enterprise to production.
>
> A monitoring platform isn't very convincing if it only runs locally, so I made sure it could be deployed for real:
>
> - Multi-stage **Docker** builds for backend, frontend, and agent (final images in the 15-40MB range)
> - Full **docker-compose** stack including PostgreSQL, Redis, and an optional observability profile (Prometheus, Grafana, Jaeger)
> - **Kubernetes manifests** — Deployment, Service, Ingress, ConfigMap, Secret — for a real cluster rollout
> - **GitHub Actions CI/CD**: lint → test → build → security scan (Trivy) → build & push images → deploy
> - A one-click cloud deployment script for spinning it up on a fresh Ubuntu VM
>
> It's now live at [your deployed URL]. Nothing teaches you more about a system than actually operating it — health checks, resource limits, and graceful shutdown all mattered a lot more once this stopped being a local demo.
>
> Next (final) post: what I'd do differently, and what's next. 🧠
>
> #devops #kubernetes #cicd #docker #buildinpublic

---

## Post 6 — Lessons Learned

**Suggested attachment:** Collage of dashboard screenshots or the architecture diagram.

> 🧠 Final post on building InfraPilot Enterprise — an infrastructure monitoring platform built with Go, React, PostgreSQL, and Redis.
>
> A few honest takeaways:
>
> 1. **Event-driven beats tightly coupled, every time.** Splitting metric ingestion from "what happens next" via an event bus made almost every later feature (alerts, audit logs, WebSocket push) trivial to add.
> 2. **Observability isn't optional, even for a monitoring tool.** Instrumenting InfraPilot itself with Prometheus/Grafana/OpenTelemetry caught real bugs before users would have.
> 3. **Polish beats scope.** It was tempting to keep bolting on features. Instead I focused the last stretch on documentation, tests, CI/CD, and deployment — the things that actually make a project feel finished rather than "in progress forever."
> 4. **Go's concurrency model is genuinely excellent for this domain.** Goroutines + channels made the agent's multi-collector design and the backend's worker pool far simpler than the equivalent in most other languages.
>
> The full source is public on GitHub, there's a live demo, and a walkthrough video — all linked below. I'm moving on to two more projects next (a CI/CD/GitOps platform and a cloud security platform) to round things out, but I'm proud of where InfraPilot ended up.
>
> 🔗 GitHub: [link]
> 🔗 Live Demo: [link]
> 🔗 Demo Video: [link]
>
> #golang #react #softwareengineering #devops #buildinpublic #portfolio
