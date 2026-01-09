'use client'

import { useEffect, useState } from 'react'
import { logs } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'

interface Log {
  id: string
  workflow_id: string
  workflow_name: string
  status: string
  message: string
  executed_at: string
}

export default function LogsPage() {
  const [logsList, setLogsList] = useState<Log[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    loadLogs()
  }, [])

  const loadLogs = async () => {
    try {
      const data = await logs.list()
      setLogsList(data || [])
    } catch (err) {
      console.error('Failed to load logs', err)
    } finally {
      setLoading(false)
    }
  }

  const filteredLogs = filter === 'all' 
    ? logsList 
    : logsList.filter(log => log.status === filter)

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Execution Logs</h1>
        <p className="text-muted-foreground">
          View execution history for all workflows
        </p>
      </div>

      <div className="flex gap-2">
        <button
          onClick={() => setFilter('all')}
          className={`px-4 py-2 rounded-md ${
            filter === 'all' ? 'bg-primary text-white' : 'bg-gray-200'
          }`}
        >
          All
        </button>
        <button
          onClick={() => setFilter('success')}
          className={`px-4 py-2 rounded-md ${
            filter === 'success' ? 'bg-green-500 text-white' : 'bg-gray-200'
          }`}
        >
          Success
        </button>
        <button
          onClick={() => setFilter('failed')}
          className={`px-4 py-2 rounded-md ${
            filter === 'failed' ? 'bg-destructive text-white' : 'bg-gray-200'
          }`}
        >
          Failed
        </button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Executions</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>Loading...</p>
          ) : filteredLogs.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-muted-foreground">No logs yet</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Workflow</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Message</TableHead>
                  <TableHead>Executed At</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredLogs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell className="font-medium">{log.workflow_name}</TableCell>
                    <TableCell>
                      <Badge variant={log.status === 'success' ? 'success' : 'destructive'}>
                        {log.status}
                      </Badge>
                    </TableCell>
                    <TableCell>{log.message}</TableCell>
                    <TableCell>{formatDate(log.executed_at)}</TableCell>
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

