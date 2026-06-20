# Nightpage Mobile App

This directory contains the source code for the native **Nightpage Mobile App** (supporting both Android and iOS), developed from scratch to extend self-hosted Nightpage instances with native mobile capabilities.

## Overview

Our mobile app allows users to seamlessly connect their self-hosted Nightpage instances to their mobile devices. 

Key benefits include:
- **Instant Connection**: Simply download the app, enter your self-hosted Nightpage URL, and log in.
- **Push Notifications**: Receive real-time alert notifications directly on your iOS or Android device.
- **On-Call On The Go**: View and manage schedules, escalations, and incident alerts directly from a native interface.

## How it Works

1. **Install the Mobile App**: Install the official hosted app from the Google Play Store or Apple App Store.
2. **Connect Instance**: On the login screen, enter the URL of your self-hosted Nightpage instance (e.g., `https://nightpage.yourdomain.com`).
3. **Log In**: Authenticate using your existing credentials, and you're ready to receive notifications and manage incidents!

## License

Unlike the core platform, which is licensed under the Apache License 2.0, the codebase in this `mobile/` directory is written from scratch and is licensed under the **GNU Affero General Public License, Version 3.0 (AGPLv3)**.

You are free to use, modify, and redistribute this mobile codebase under the terms of the AGPLv3. Please ensure that if you run a modified version of the mobile application as a service, the source code of your modifications is made available to the public.

## Hosting & Paid Features

To cover ongoing server infrastructure, push notification gateway services, and maintenance costs associated with hosting the centralized mobile app services, **paid features and subscriptions** will be introduced in the future for the hosted mobile app service. Connecting to self-hosted instances using the basic features will remain free and open.

## Contributions

Contributions are welcome! Please make sure any contributions to this folder are compatible with the AGPL-3.0 license.
