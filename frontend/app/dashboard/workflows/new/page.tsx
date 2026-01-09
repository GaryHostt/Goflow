'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { workflows } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export default function NewWorkflowPage() {
  const router = useRouter()
  const [name, setName] = useState('')
  const [triggerType, setTriggerType] = useState('webhook')
  const [actionType, setActionType] = useState('slack_message')
  const [config, setConfig] = useState({
    slack_message: '',
    discord_message: '',
    city: '',
    interval: 10,
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      // Build config based on action type
      let configJson: any = {}
      
      if (actionType === 'slack_message') {
        configJson.slack_message = config.slack_message || 'Hello from iPaaS!'
      } else if (actionType === 'discord_post') {
        configJson.discord_message = config.discord_message || 'Hello from iPaaS!'
      } else if (actionType === 'weather_check') {
        configJson.city = config.city || 'London'
      }

      if (triggerType === 'schedule') {
        configJson.interval = config.interval
      }

      await workflows.create({
        name,
        trigger_type: triggerType,
        action_type: actionType,
        config_json: JSON.stringify(configJson),
      })

      router.push('/dashboard/workflows')
    } catch (err) {
      setError('Failed to create workflow. Make sure you have connected the required service.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Create Workflow</h1>
        <p className="text-muted-foreground">
          Set up a new integration workflow
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Workflow Configuration</CardTitle>
            <CardDescription>
              Configure your trigger and action
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {error && (
              <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <label className="text-sm font-medium">Workflow Name</label>
              <Input
                placeholder="My Integration Workflow"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Trigger Type</label>
              <select
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                value={triggerType}
                onChange={(e) => setTriggerType(e.target.value)}
              >
                <option value="webhook">Webhook</option>
                <option value="schedule">Schedule (Polling)</option>
              </select>
              {triggerType === 'webhook' && (
                <p className="text-xs text-muted-foreground">
                  A unique webhook URL will be generated after creation
                </p>
              )}
            </div>

            {triggerType === 'schedule' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Interval (minutes)</label>
                <Input
                  type="number"
                  min="1"
                  value={config.interval}
                  onChange={(e) => setConfig({ ...config, interval: parseInt(e.target.value) })}
                />
              </div>
            )}

            <div className="space-y-2">
              <label className="text-sm font-medium">Action Type</label>
              <select
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                value={actionType}
                onChange={(e) => setActionType(e.target.value)}
              >
                <option value="slack_message">Send Slack Message</option>
                <option value="discord_post">Send Discord Message</option>
                <option value="weather_check">Check Weather</option>
              </select>
            </div>

            {actionType === 'slack_message' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Slack Message</label>
                <Input
                  placeholder="Hello from iPaaS!"
                  value={config.slack_message}
                  onChange={(e) => setConfig({ ...config, slack_message: e.target.value })}
                />
              </div>
            )}

            {actionType === 'discord_post' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Discord Message</label>
                <Input
                  placeholder="Hello from iPaaS!"
                  value={config.discord_message}
                  onChange={(e) => setConfig({ ...config, discord_message: e.target.value })}
                />
              </div>
            )}

            {actionType === 'weather_check' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">City</label>
                <Input
                  placeholder="London"
                  value={config.city}
                  onChange={(e) => setConfig({ ...config, city: e.target.value })}
                />
              </div>
            )}

            <div className="flex gap-4">
              <Button type="submit" disabled={loading}>
                {loading ? 'Creating...' : 'Create Workflow'}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push('/dashboard/workflows')}
              >
                Cancel
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>
    </div>
  )
}

