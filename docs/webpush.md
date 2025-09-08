# Web Push Notifications

GoAlert can send browser push notifications using the Web Push protocol. This feature uses VAPID keys and requires HTTPS.

## Generating VAPID Keys

Use the [`webpush-go`](https://github.com/SherClockHolmes/webpush-go) tool to generate a key pair:

```
go run github.com/SherClockHolmes/webpush-go@latest -gen
```

Copy the printed public and private keys into the configuration as `WebPush.VAPIDPublicKey` and `WebPush.VAPIDPrivateKey`.

## Configuration

Set the following in the configuration file or environment:

```
WebPush.Enable = true
WebPush.VAPIDPublicKey = <public key>
WebPush.VAPIDPrivateKey = <private key>
WebPush.TTL = 30
```

The TTL defines how long (in seconds) push messages should be retained by the push service.

## Browser Requirements

Web Push requires HTTPS and a compatible browser. Most modern browsers support the Push API, but users must grant permission to receive notifications.

## Progressive Web App

GoAlert includes a service worker and manifest so it can be installed as a Progressive Web App. Most browsers provide an "Add to Home Screen" option that installs the app for quick access without opening the browser first.
