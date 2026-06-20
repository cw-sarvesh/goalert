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
