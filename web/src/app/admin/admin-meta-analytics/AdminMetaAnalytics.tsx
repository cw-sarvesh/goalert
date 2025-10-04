import React, { useCallback, useMemo, useState } from 'react'
import {
  Alert,
  Button,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Grid,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import { DateTime } from 'luxon'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import AutoSizer from 'react-virtualized-auto-sizer'
import { useTheme } from '@mui/material/styles'

interface MetaPriorityCount {
  value: string
  count: number
}

interface MetaAckLevelCount {
  escalationLevel: number | null
  count: number
}

interface MetaAnalytics {
  priorityCounts: MetaPriorityCount[]
  acknowledgedByLevel: MetaAckLevelCount[]
}

const formatInput = (dt: DateTime): string =>
  dt.toISO({ suppressMilliseconds: true, includeOffset: false }) || ''

export default function AdminMetaAnalytics(): React.JSX.Element {
  const theme = useTheme()
  const now = useMemo(() => DateTime.now(), [])
  const [serviceID, setServiceID] = useState('')
  const [metaKey, setMetaKey] = useState('alerts/priority')
  const [startValue, setStartValue] = useState(formatInput(now.minus({ days: 7 })))
  const [endValue, setEndValue] = useState(formatInput(now))
  const [analytics, setAnalytics] = useState<MetaAnalytics | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleFetch = useCallback(async () => {
    setError(null)

    const startDT = DateTime.fromISO(startValue)
    if (!startDT.isValid) {
      setError('Invalid start time. Please provide a valid value.')
      return
    }

    const endDT = DateTime.fromISO(endValue)
    if (!endDT.isValid) {
      setError('Invalid end time. Please provide a valid value.')
      return
    }

    if (!serviceID.trim()) {
      setError('Service ID is required.')
      return
    }

    if (!metaKey.trim()) {
      setError('Metadata key is required.')
      return
    }

    if (startDT.toMillis() >= endDT.toMillis()) {
      setError('Start time must be before end time.')
      return
    }

    const params = new URLSearchParams({
      serviceID: serviceID.trim(),
      metaKey: metaKey.trim(),
      start: startDT.toUTC().toISO({ suppressMilliseconds: true }) || '',
      end: endDT.toUTC().toISO({ suppressMilliseconds: true }) || '',
    })

    setLoading(true)
    try {
      const res = await fetch(`/admin/analytics/service-meta?${params.toString()}`)
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Failed to fetch analytics data')
      }
      const data = (await res.json()) as MetaAnalytics
      setAnalytics(data)
    } catch (err) {
      setAnalytics(null)
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }, [startValue, endValue, serviceID, metaKey])

  const renderPriorityChart = (data: MetaPriorityCount[]): React.JSX.Element => {
    if (!data.length) {
      return <Typography>No alerts matched the selected metadata.</Typography>
    }

    const chartData = data.map(({ value, count }) => ({ value, count }))

    return (
      <AutoSizer>
        {({ width, height }) => (
          <BarChart
            width={width}
            height={height}
            data={chartData}
            margin={{ top: 40, right: 24, bottom: 32 }}
          >
            <CartesianGrid strokeDasharray='3 3' />
            <XAxis dataKey='value' stroke={theme.palette.text.secondary} />
            <YAxis stroke={theme.palette.text.secondary} allowDecimals={false} />
            <Tooltip />
            <Legend />
            <Bar
              dataKey='count'
              name='Alert Count'
              fill={theme.palette.primary.main}
            />
          </BarChart>
        )}
      </AutoSizer>
    )
  }

  const renderAckChart = (data: MetaAckLevelCount[]): React.JSX.Element => {
    if (!data.length) {
      return (
        <Typography>
          No acknowledgements recorded for alerts with the selected metadata.
        </Typography>
      )
    }

    const chartData = data.map(({ escalationLevel, count }) => ({
      level:
        escalationLevel === null || escalationLevel === undefined || escalationLevel < 0
          ? 'Unspecified'
          : `Level ${escalationLevel + 1}`,
      count,
    }))

    return (
      <AutoSizer>
        {({ width, height }) => (
          <BarChart
            width={width}
            height={height}
            data={chartData}
            margin={{ top: 40, right: 24, bottom: 32 }}
          >
            <CartesianGrid strokeDasharray='3 3' />
            <XAxis dataKey='level' stroke={theme.palette.text.secondary} />
            <YAxis stroke={theme.palette.text.secondary} allowDecimals={false} />
            <Tooltip />
            <Legend />
            <Bar
              dataKey='count'
              name='Acknowledgements'
              fill={theme.palette.secondary.main}
            />
          </BarChart>
        )}
      </AutoSizer>
    )
  }

  const totalAlerts = analytics?.priorityCounts?.reduce(
    (acc, cur) => acc + cur.count,
    0,
  )
  const totalAcks = analytics?.acknowledgedByLevel?.reduce(
    (acc, cur) => acc + cur.count,
    0,
  )

  return (
    <Grid container spacing={2}>
      <Grid item xs={12}>
        <Card>
          <CardHeader
            title='Alert Priority Analytics'
            subheader='Visualise alert metadata distribution and acknowledgement levels'
          />
          <CardContent>
            <form
              onSubmit={(evt) => {
                evt.preventDefault()
                handleFetch()
              }}
            >
              <Stack
                direction={{ xs: 'column', md: 'row' }}
                spacing={2}
                alignItems={{ xs: 'stretch', md: 'flex-end' }}
              >
                <TextField
                  label='Service ID'
                  value={serviceID}
                  onChange={(e) => setServiceID(e.target.value)}
                  required
                  fullWidth
                />
                <TextField
                  label='Metadata Key'
                  value={metaKey}
                  onChange={(e) => setMetaKey(e.target.value)}
                  required
                  fullWidth
                />
                <TextField
                  label='Start'
                  type='datetime-local'
                  value={startValue}
                  onChange={(e) => setStartValue(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                  required
                />
                <TextField
                  label='End'
                  type='datetime-local'
                  value={endValue}
                  onChange={(e) => setEndValue(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                  required
                />
                <Button
                  type='submit'
                  variant='contained'
                  disabled={loading}
                  sx={{ minWidth: 160 }}
                >
                  {loading ? <CircularProgress size={24} /> : 'Load Analytics'}
                </Button>
              </Stack>
            </form>
            {error && (
              <Alert severity='error' sx={{ mt: 2 }}>
                {error}
              </Alert>
            )}
          </CardContent>
        </Card>
      </Grid>

      {analytics && !loading && (
        <>
          <Grid item xs={12} md={6} sx={{ height: 420 }}>
            <Card sx={{ height: '100%' }}>
              <CardHeader title='Alerts by Metadata Value' />
              <CardContent sx={{ height: '100%' }}>
                {renderPriorityChart(analytics.priorityCounts || [])}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={6} sx={{ height: 420 }}>
            <Card sx={{ height: '100%' }}>
              <CardHeader title='Acknowledgements by Escalation Level' />
              <CardContent sx={{ height: '100%' }}>
                {renderAckChart(analytics.acknowledgedByLevel || [])}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12}>
            <Card>
              <CardHeader title='Summary' />
              <CardContent>
                <Stack spacing={1}>
                  <Typography>Total Alerts: {totalAlerts ?? 0}</Typography>
                  <Typography>Total Acknowledgements: {totalAcks ?? 0}</Typography>
                </Stack>
              </CardContent>
            </Card>
          </Grid>
        </>
      )}
    </Grid>
  )
}
