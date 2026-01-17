'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { workflows } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { WorkflowFlowDiagram } from '@/components/WorkflowFlowDiagram'
import { Input } from '@/components/ui/input'

export default function NewWorkflowPage() {
  const router = useRouter()
  const [name, setName] = useState('')
  const [triggerType, setTriggerType] = useState('webhook')
  const [actionType, setActionType] = useState('slack_message')
  const [webhookPath, setWebhookPath] = useState('/integration/my-endpoint')
  const [webhookPayload, setWebhookPayload] = useState('{\n  "user": {\n    "name": ""\n  }\n}')
  const [config, setConfig] = useState<any>({
    slack_message: '',
    discord_message: '',
    twilio_to: '',
    twilio_message: '',
    city: '',
    interval: 10,
    testing_response_json: '{"message": "Test response", "status": "success"}',
    testing_status_code: 200,
    testing_delay: 0,
    news_query: '',
    news_country: 'us',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async () => {
    if (!name) {
      setError('Please enter a workflow name')
      return
    }

    setError('')
    setLoading(true)

    try {
      // Build config based on action type
      let configJson: any = { ...config }
      
      if (triggerType === 'webhook') {
        configJson.webhook_path = webhookPath
        configJson.webhook_payload = webhookPayload
      } else if (triggerType === 'schedule') {
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
    <div className="container mx-auto py-6 space-y-6">
      {/* Header with Name and Actions */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold">Create Workflow</h1>
          <p className="text-muted-foreground">
            Click on elements below to configure your integration
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => router.push('/dashboard/workflows')}
          >
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={loading || !name}>
            {loading ? 'Creating...' : 'Create Workflow'}
          </Button>
        </div>
      </div>

      {/* Workflow Name Input */}
      <Card>
        <CardContent className="pt-6">
          <div className="max-w-md">
            <label className="text-sm font-medium mb-2 block">Workflow Name</label>
            <Input
              placeholder="My Integration Workflow"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          {error && (
            <div className="mt-4 rounded-md bg-destructive/15 p-3 text-sm text-destructive">
              {error}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Visual Flow - Main Focus */}
      <Card>
        <CardHeader>
          <CardTitle>Workflow Builder</CardTitle>
          <CardDescription>
            Click on the trigger and action to configure them
          </CardDescription>
        </CardHeader>
        <CardContent className="min-h-[500px]">
          <WorkflowFlowDiagram 
            triggerType={triggerType}
            actionType={actionType}
            config={config}
            webhookPath={webhookPath}
            webhookPayload={webhookPayload}
            onTriggerTypeChange={setTriggerType}
            onActionTypeChange={setActionType}
            onWebhookPathChange={setWebhookPath}
            onWebhookPayloadChange={setWebhookPayload}
            onConfigChange={setConfig}
          />
        </CardContent>
      </Card>
    </div>
  )
}
