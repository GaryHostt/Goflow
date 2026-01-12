'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { workflows } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Loader2, AlertCircle, Inbox } from 'lucide-react'

interface Workflow {
  id: string
  name: string
  trigger_type: string
  action_type: string
  is_active: boolean
  created_at: string
}

export default function WorkflowsPage() {
  const [workflowsList, setWorkflowsList] = useState<Workflow[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>('')

  const loadWorkflows = async () => {
    try {
      setError('')
      const data = await workflows.list()
      setWorkflowsList(data || [])
    } catch (err: any) {
      console.error('Failed to load workflows', err)
      setError(err.message || 'Failed to load workflows. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadWorkflows()
  }, [])

  const handleToggle = async (id: string) => {
    try {
      await workflows.toggle(id)
      await loadWorkflows()
    } catch (err: any) {
      console.error('Failed to toggle workflow', err)
      setError(err.message || 'Failed to toggle workflow')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this workflow?')) return
    
    try {
      await workflows.delete(id)
      await loadWorkflows()
    } catch (err: any) {
      console.error('Failed to delete workflow', err)
      setError(err.message || 'Failed to delete workflow')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Workflows</h1>
          <p className="text-muted-foreground">
            Manage your integration workflows
          </p>
        </div>
        <Link href="/dashboard/workflows/new">
          <Button>Create Workflow</Button>
        </Link>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Your Workflows</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex flex-col items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-primary mb-4" />
              <p className="text-muted-foreground">Loading workflows...</p>
            </div>
          ) : workflowsList.length === 0 ? (
            <div className="text-center py-12">
              <Inbox className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">No workflows yet</h3>
              <p className="text-muted-foreground mb-6">
                Get started by creating your first integration workflow
              </p>
              <Link href="/dashboard/workflows/new">
                <Button>Create Your First Workflow</Button>
              </Link>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Trigger</TableHead>
                  <TableHead>Action</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {workflowsList.map((workflow) => (
                  <TableRow key={workflow.id}>
                    <TableCell className="font-medium">{workflow.name}</TableCell>
                    <TableCell className="capitalize">{workflow.trigger_type}</TableCell>
                    <TableCell className="capitalize">
                      {workflow.action_type.replace('_', ' ')}
                    </TableCell>
                    <TableCell>
                      <Badge variant={workflow.is_active ? 'success' : 'secondary'}>
                        {workflow.is_active ? 'Active' : 'Inactive'}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex gap-2">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleToggle(workflow.id)}
                        >
                          {workflow.is_active ? 'Disable' : 'Enable'}
                        </Button>
                        <Button
                          size="sm"
                          variant="destructive"
                          onClick={() => handleDelete(workflow.id)}
                        >
                          Delete
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

