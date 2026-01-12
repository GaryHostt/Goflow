"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Shield, Zap, Database, Lock, DollarSign, Plus, Trash2, Network, Webhook, Layers } from "lucide-react";
import { api } from "@/lib/api";

interface Workflow {
  id: string;
  name: string;
  action_type: string;
  is_active: boolean;
}

interface KongService {
  id: string;
  name: string;
  url: string;
  created_at: number;
}

interface UseCase {
  id: string;
  name: string;
  description: string;
  icon: any;
  benefits: string[];
}

export default function APIManagementPage() {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [kongServices, setKongServices] = useState<KongService[]>([]);
  const [selectedWorkflow, setSelectedWorkflow] = useState<string>("");
  const [selectedUseCase, setSelectedUseCase] = useState<string>("");
  const [serviceName, setServiceName] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Define 5 Kong use cases
  const useCases: UseCase[] = [
    {
      id: "protocol_bridge",
      name: "Protocol Bridge",
      description: "Convert SOAP/Legacy APIs to modern REST endpoints",
      icon: Layers,
      benefits: [
        "Hide legacy system complexity",
        "Modern JSON API interface",
        "Automatic protocol conversion"
      ]
    },
    {
      id: "webhook_handler",
      name: "Webhook Handler",
      description: "Rate-limited webhook processing with flow control",
      icon: Webhook,
      benefits: [
        "Prevent webhook flooding",
        "100 req/sec rate limiting",
        "Automatic retry handling"
      ]
    },
    {
      id: "aggregator",
      name: "Smart API Aggregator",
      description: "Combine multiple APIs with caching",
      icon: Database,
      benefits: [
        "Reduce client-server chattiness",
        "5-minute response caching",
        "Parallel API orchestration"
      ]
    },
    {
      id: "auth_overlay",
      name: "Federated Security",
      description: "Add OAuth2/API Key auth to any workflow",
      icon: Lock,
      benefits: [
        "Centralized authentication",
        "API key or OAuth2 support",
        "Trust header injection"
      ]
    },
    {
      id: "monetization",
      name: "Usage-Based Billing",
      description: "Track API usage for monetization",
      icon: DollarSign,
      benefits: [
        "Request counting & limits",
        "ELK usage analytics",
        "Pay-per-use ready"
      ]
    }
  ];

  useEffect(() => {
    loadWorkflows();
    loadKongServices();
  }, []);

  const loadWorkflows = async () => {
    try {
      const response = await api.get("/api/workflows");
      setWorkflows(response.data.data || []);
    } catch (err: any) {
      setError("Failed to load workflows");
    }
  };

  const loadKongServices = async () => {
    try {
      const response = await api.get("/api/kong/services");
      const services = response.data.data?.data || [];
      setKongServices(services);
    } catch (err: any) {
      // Kong might not be running, that's okay
      console.log("Kong not available:", err);
    }
  };

  const createKongTemplate = async () => {
    if (!selectedWorkflow || !selectedUseCase || !serviceName) {
      setError("Please select a workflow, use case, and provide a service name");
      return;
    }

    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await api.post("/api/kong/templates", {
        workflow_id: selectedWorkflow,
        use_case: selectedUseCase,
      });

      setSuccess(`Kong service "${serviceName}" created successfully with ${selectedUseCase} template!`);
      loadKongServices();
      setServiceName("");
      setSelectedWorkflow("");
      setSelectedUseCase("");
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create Kong service");
    } finally {
      setLoading(false);
    }
  };

  const deleteKongService = async (serviceId: string) => {
    if (!confirm("Are you sure you want to delete this Kong service?")) {
      return;
    }

    try {
      await api.delete(`/api/kong/services/${serviceId}`);
      setSuccess("Kong service deleted successfully");
      loadKongServices();
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to delete Kong service");
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">API Management</h1>
          <p className="text-gray-500 mt-1">
            Add enterprise-grade API Gateway features to your workflows using Kong
          </p>
        </div>
        <Badge variant="secondary" className="h-fit">
          <Network className="w-4 h-4 mr-1" />
          Powered by Kong Gateway
        </Badge>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {success && (
        <Alert className="border-green-500 bg-green-50">
          <AlertDescription className="text-green-800">{success}</AlertDescription>
        </Alert>
      )}

      {/* Use Case Templates */}
      <Card>
        <CardHeader>
          <CardTitle>5 Enterprise Use Cases</CardTitle>
          <CardDescription>
            Click a use case to apply Kong Gateway patterns to your workflows
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {useCases.map((useCase) => {
              const Icon = useCase.icon;
              return (
                <div
                  key={useCase.id}
                  className={`border rounded-lg p-4 cursor-pointer transition-all hover:shadow-md ${
                    selectedUseCase === useCase.id
                      ? "border-blue-500 bg-blue-50"
                      : "border-gray-200"
                  }`}
                  onClick={() => setSelectedUseCase(useCase.id)}
                >
                  <div className="flex items-start space-x-3">
                    <div className="bg-blue-100 p-2 rounded">
                      <Icon className="w-5 h-5 text-blue-600" />
                    </div>
                    <div className="flex-1">
                      <h3 className="font-semibold text-sm">{useCase.name}</h3>
                      <p className="text-xs text-gray-600 mt-1">{useCase.description}</p>
                      <ul className="mt-2 space-y-1">
                        {useCase.benefits.map((benefit, index) => (
                          <li key={index} className="text-xs text-gray-500 flex items-start">
                            <span className="text-green-500 mr-1">✓</span>
                            {benefit}
                          </li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>

      {/* Create Kong Service */}
      <Card>
        <CardHeader>
          <CardTitle>Create Kong Service</CardTitle>
          <CardDescription>
            Apply a Kong template to one of your existing workflows
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Service Name</label>
            <Input
              placeholder="e.g., legacy-soap-bridge"
              value={serviceName}
              onChange={(e) => setServiceName(e.target.value)}
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Select Workflow</label>
            <Select value={selectedWorkflow} onValueChange={setSelectedWorkflow}>
              <SelectTrigger>
                <SelectValue placeholder="Choose a workflow to expose via Kong" />
              </SelectTrigger>
              <SelectContent>
                {workflows
                  .filter((wf) => wf.is_active)
                  .map((workflow) => (
                    <SelectItem key={workflow.id} value={workflow.id}>
                      {workflow.name} ({workflow.action_type})
                    </SelectItem>
                  ))}
              </SelectContent>
            </Select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Select Use Case Template</label>
            <Select value={selectedUseCase} onValueChange={setSelectedUseCase}>
              <SelectTrigger>
                <SelectValue placeholder="Choose a Kong template" />
              </SelectTrigger>
              <SelectContent>
                {useCases.map((useCase) => (
                  <SelectItem key={useCase.id} value={useCase.id}>
                    {useCase.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <Button
            onClick={createKongTemplate}
            disabled={loading || !selectedWorkflow || !selectedUseCase || !serviceName}
            className="w-full"
          >
            <Plus className="w-4 h-4 mr-2" />
            {loading ? "Creating..." : "Create Kong Service"}
          </Button>
        </CardContent>
      </Card>

      {/* Existing Kong Services */}
      <Card>
        <CardHeader>
          <CardTitle>Active Kong Services</CardTitle>
          <CardDescription>
            Manage your Kong Gateway services and routes
          </CardDescription>
        </CardHeader>
        <CardContent>
          {kongServices.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <Network className="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p>No Kong services yet. Create one above to get started!</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Upstream URL</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {kongServices.map((service) => (
                  <TableRow key={service.id}>
                    <TableCell className="font-medium">{service.name}</TableCell>
                    <TableCell className="text-sm text-gray-600">{service.url}</TableCell>
                    <TableCell className="text-sm text-gray-600">
                      {new Date(service.created_at * 1000).toLocaleDateString()}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => deleteKongService(service.id)}
                      >
                        <Trash2 className="w-4 h-4 text-red-500" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Kong Manager Link */}
      <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-blue-200">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-semibold text-lg mb-1">Kong Manager (GUI)</h3>
              <p className="text-sm text-gray-600">
                Access the full Kong Gateway admin interface for advanced configuration
              </p>
            </div>
            <Button asChild variant="outline">
              <a href="http://localhost:8002" target="_blank" rel="noopener noreferrer">
                <Shield className="w-4 h-4 mr-2" />
                Open Kong Manager
              </a>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

