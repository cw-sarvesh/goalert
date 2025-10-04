export {}

declare global {
  interface Window {
    goalertPush?: {
      ensureSubscribed: () => Promise<PushSubscription | void>
      requestPermissionAndSubscribe: () => Promise<PushSubscription | void>
      getPermission: () => NotificationPermission
    }
  }
}
