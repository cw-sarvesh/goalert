import React, { useCallback, useEffect, useMemo, useState } from 'react'
import Dialog from '@mui/material/Dialog'
import DialogTitle from '@mui/material/DialogTitle'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import DialogContentText from '@mui/material/DialogContentText'
import Button from '@mui/material/Button'
import Alert from '@mui/material/Alert'

const hasPushSupport = (): boolean => {
  return (
    typeof window !== 'undefined' &&
    'Notification' in window &&
    'serviceWorker' in navigator &&
    'PushManager' in window
  )
}

const getPermission = (): NotificationPermission => {
  if (typeof window === 'undefined') return 'denied'
  const fromGlobal = window.goalertPush?.getPermission?.()
  if (fromGlobal) return fromGlobal
  if ('Notification' in window) return Notification.permission
  return 'denied'
}

export default function NotificationPermissionPrompt(): React.JSX.Element | null {
  const supported = useMemo(hasPushSupport, [])
  const [permission, setPermission] = useState<NotificationPermission>(getPermission)
  const [dismissed, setDismissed] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const refreshPermission = useCallback(() => {
    setPermission(getPermission())
  }, [])

  useEffect(() => {
    if (!supported) return

    refreshPermission()

    const handleVisibility = () => refreshPermission()
    window.addEventListener('focus', handleVisibility)
    document.addEventListener('visibilitychange', handleVisibility)

    let permissionStatus: PermissionStatus | null = null
    if ('permissions' in navigator && navigator.permissions.query) {
      navigator.permissions
        .query({ name: 'notifications' as PermissionName })
        .then((status) => {
          permissionStatus = status
          status.onchange = refreshPermission
        })
        .catch(() => {
          /* ignore */
        })
    }

    return () => {
      window.removeEventListener('focus', handleVisibility)
      document.removeEventListener('visibilitychange', handleVisibility)
      if (permissionStatus) permissionStatus.onchange = null
    }
  }, [refreshPermission, supported])

  useEffect(() => {
    if (permission === 'granted') {
      setDismissed(false)
      setError(null)
    }
  }, [permission])

  const handleEnable = () => {
    setError(null)
    if (!window.goalertPush?.requestPermissionAndSubscribe) {
      setError('Push setup is not ready yet. Please try again in a moment.')
      return
    }
    setLoading(true)
    window.goalertPush
      .requestPermissionAndSubscribe()
      .then(() => {
        refreshPermission()
        setDismissed(false)
      })
      .catch((err: unknown) => {
        if (err instanceof Error) {
          setError(err.message)
        } else if (typeof err === 'string') {
          setError(err)
        } else {
          setError('Unable to enable notifications.')
        }
      })
      .finally(() => setLoading(false))
  }

  const handleLater = () => {
    setDismissed(true)
  }

  if (!supported) return null

  const open = permission === 'default' && !dismissed

  return (
    <Dialog open={open} onClose={handleLater} aria-labelledby='notification-permission-title'>
      <DialogTitle id='notification-permission-title'>Enable Notifications</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Allow browser notifications to receive real-time alerts on this device.
        </DialogContentText>
        {error ? <Alert severity='error' sx={{ mt: 2 }}>{error}</Alert> : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleLater} color='inherit'>
          Not now
        </Button>
        <Button onClick={handleEnable} variant='contained' disabled={loading} autoFocus>
          Enable notifications
        </Button>
      </DialogActions>
    </Dialog>
  )
}
