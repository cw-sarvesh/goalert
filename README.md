# NOTE
This project, **Nightpage**, is a forked and modified version of [GoAlert](https://github.com/target/goalert/). We have implemented several modifications and enhancements, and have uploaded the revised version to our repository under the new name.

## Features
1. Progressive Web App (PWA) support
2. Web push notifications
3. Alert segregation: Phone calls for critical (P1) alerts, and push notifications for P2 and other alerts
4. Option to manually disable SMS notifications
5. Option to manually disable call notifications
6. **Mobile App Support (Android & iOS)**: Connect your self-hosted Nightpage instance directly to our mobile applications (codebase located under the [`mobile/`](./mobile) directory).


# Nightpage

Nightpage provides on-call scheduling, automated escalations and notifications (like SMS or voice calls) to automatically engage the right person, the right way, and at the right time.

![main-screen-updated](https://user-images.githubusercontent.com/595010/189744659-66ee6aed-b7b6-4625-a2ac-1f8ad3c1ea4f.png)

## Installation

Nightpage is built and distributed as a single binary.

See our [Getting Started Guide](./docs/getting-started.md) for running Nightpage in a production environment.

### Quick Start

*(Note: The demo container below is provided by the upstream GoAlert project)*

```bash
docker run -it --rm -p 8081:8081 goalert/demo
```

The application will be running at [localhost:8081](http://localhost:8081). You can log in with `admin`/`admin123`.

If you're using the demo container for integration testing:

- A non-admin user is available as `user`/`user1234`.
- You can specify the ENV variable `SKIP_SEED=1` to skip the initial seed data step.
- You can get a session token via `curl -XPOST -H 'Referer: http://localhost:8081' -d 'username=admin&password=admin123' 'http://localhost:8081/api/v2/identity/providers/basic?noRedirect=1'`.

## Building Docker images and testing locally

You can build GoAlert Docker images yourself and run them for local testing, or run GoAlert on your host with a local or containerized database.

### Build the Docker images

From the repository root, with `make` and `docker` available:

| Image | Command | Description |
|-------|---------|-------------|
| **GoAlert** (app only) | `make container-goalert` | Builds the GoAlert binary for Linux (amd64/arm/arm64) and produces `goalert/goalert:<version>`. Use this image when you have PostgreSQL running elsewhere. |
| **Demo** (app + PostgreSQL) | `make container-demo` | Builds the demo image that includes GoAlert and an embedded PostgreSQL. Produces `goalert/demo:<version>`. Easiest for a single-command local run. |

The image tag is derived from git (e.g. `v1.2.3` or `dev`). To push images to a registry, use `make container-goalert PUSH=1` (or `container-demo`).

### Test with Docker

**Option 1 — Demo image (all-in-one, good for a quick try):**

```bash
# Build and run the demo image (GoAlert + PostgreSQL in one container)
make container-demo
docker run -it --rm -p 8081:8081 goalert/demo:dev
```

Open [http://localhost:8081](http://localhost:8081) and log in with `admin` / `admin123`. Replace `dev` with your current git-based tag if different.

**Option 2 — GoAlert image with a database on the host or in another container:**

```bash
# Start PostgreSQL (e.g. in Docker)
make postgres

# Build the GoAlert image
make container-goalert

# Run GoAlert, connecting to PostgreSQL on the host
docker run -it --rm -p 8081:8081 \
  -e GOALERT_DB_URL=postgres://goalert@host.docker.internal:5432/goalert \
  -e GOALERT_DATA_ENCRYPTION_KEY=your-secret-key \
  -e GOALERT_PUBLIC_URL=http://localhost:8081 \
  goalert/goalert:dev
```

On Linux, use `-e GOALERT_DB_URL=postgres://goalert@<host-ip>:5432/goalert` and ensure PostgreSQL accepts connections from the Docker network, or run Postgres in Docker and link the containers.

### Test on the host (without running GoAlert in Docker)

Run GoAlert as a local binary with hot-reload; PostgreSQL can still run in Docker:

```bash
# Start PostgreSQL in Docker (if not already running)
make postgres

# Start GoAlert and dev tooling (UI + API + Prometheus, Mailpit, etc.)
make start
```

- GoAlert UI and API: [http://localhost:3030](http://localhost:3030)
- Default login: `admin` / `admin123`
- Code and config changes are picked up without restarting the server.

To run only the database in Docker and start GoAlert yourself:

```bash
make postgres
./bin/goalert --db-url postgres://goalert@localhost:5432/goalert \
  --data-encryption-key dev-key --public-url http://localhost:3030
```

For more detail on production-style deployment and configuration, see the [Getting Started Guide](./docs/getting-started.md). For full development setup (tools, DB, resetdb, etc.), see [Development Setup](./docs/development-setup.md).

## Contributing (Local Development)

If you'd like to contribute to Nightpage, please see our [Contributing Guidelines](./CONTRIBUTING.md) and the [Development Setup Guide](./docs/development-setup.md).

Please also see our [Code of Conduct](./CODE_OF_CONDUCT.md).

For most purposes, you can use `make start` from the root of this repo to start a development server.

- It will be running at `http://localhost:3030`
- Default login is `admin`/`admin123`
- Changes you make locally (UI and backend) should be reflected in the running server within a few seconds (no need to restart the server).

## Upstream Project

Nightpage is a fork of GoAlert. For help, issues, or discussions regarding the original GoAlert platform:
- Slack: The `#goalert` channel on [gophers.slack.com](https://gophers.slack.com/messages/goalert/) (requires an invitation from https://invite.slack.golangbridge.org/)
- Bug Reports / Feature Requests: [target/goalert Issues](https://github.com/target/goalert/issues)

## License

This repository is dual-licensed to accommodate both the core platform and the custom mobile application features:

- **Core Web Platform & Server**: The server and web platform (originally derived from GoAlert) are licensed under the [Apache License, Version 2.0](./LICENSE.md).
- **Mobile Application**: The source code in the [`mobile/`](./mobile) directory is developed from scratch and is licensed under the [GNU Affero General Public License (AGPL) v3.0](https://www.gnu.org/licenses/agpl-3.0.html) (see [`mobile/README.md`](./mobile/README.md) for more details).
