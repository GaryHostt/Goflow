'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { workflows } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

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

  const loadWorkflows = async () => {
    try {
      const data = await workflows.list()
      setWorkflowsList(data || [])
    } catch (err) {
      console.error('Failed to load workflows', err)
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
    } catch (err) {
      console.error('Failed to toggle workflow', err)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this workflow?')) return
    
    try {
      await workflows.delete(id)
      await loadWorkflows()
    } catch (err) {
      console.error('Failed to delete workflow', err)
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

      <Card>
        <CardHeader>
          <CardTitle>Your Workflows</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>Loading...</p>
          ) : workflowsList.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground mb-4">No workflows yet</p>
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

