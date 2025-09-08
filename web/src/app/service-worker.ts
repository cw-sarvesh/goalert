const prefix = self.location.pathname.replace(
  /\/static\/service-worker\.js$/,
  '',
)
const CACHE_NAME = 'goalert-cache-v1'
const PRECACHE_URLS = [
  prefix + '/',
  prefix + '/static/app.js',
  prefix + '/static/app.css',
]

self.addEventListener('install', (event: ExtendableEvent) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(PRECACHE_URLS)),
  )
})

self.addEventListener('fetch', (event: FetchEvent) => {
  event.respondWith(
    caches.match(event.request).then((resp) => resp || fetch(event.request)),
  )
})

self.addEventListener('push', (event: PushEvent) => {
  const data = event.data?.json() || {}
  const title = data.title || 'GoAlert'
  const body = data.body || ''
  const url = data.url || prefix + '/'
  event.waitUntil(
    self.registration.showNotification(title, {
      body,
      data: { url },
    }),
  )

  event.waitUntil(
    self.clients
      .matchAll({ type: 'window', includeUncontrolled: true })
      .then((clients) => {
        for (const client of clients) {
          client.postMessage(data)
        }
      }),
  )
})

self.addEventListener('notificationclick', (event: NotificationEvent) => {
  event.notification.close()
  const url = event.notification.data?.url || prefix + '/'
  event.waitUntil(
    self.clients
      .matchAll({ type: 'window', includeUncontrolled: true })
      .then((clientList) => {
        for (const client of clientList) {
          if ('focus' in client) return client.focus()
        }
        if (self.clients.openWindow) return self.clients.openWindow(url)
      }),
  )
})
