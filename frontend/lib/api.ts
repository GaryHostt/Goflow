// API base URL - update this to match your backend
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Get auth token from localStorage
const getToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
};

// Set auth token in localStorage
export const setToken = (token: string): void => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('token', token);
  }
};

// Clear auth token
export const clearToken = (): void => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
  }
};

// Check if user is authenticated
export const isAuthenticated = (): boolean => {
  return getToken() !== null;
};

// API client with JWT injection
// TODO: MULTI-TENANT - Add tenant_id to headers
async function apiClient(endpoint: string, options: RequestInit = {}) {
  const token = getToken();
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    clearToken();
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
  }

  return response;
}

// Auth API
export const auth = {
  register: async (email: string, password: string) => {
    const response = await apiClient('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    return response.json();
  },
  
  login: async (email: string, password: string) => {
    const response = await apiClient('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    return response.json();
  },
};

// Credentials API
export const credentials = {
  create: async (serviceName: string, apiKey: string) => {
    const response = await apiClient('/credentials', {
      method: 'POST',
      body: JSON.stringify({ service_name: serviceName, api_key: apiKey }),
    });
    return response.json();
  },
  
  list: async () => {
    const response = await apiClient('/credentials');
    return response.json();
  },
};

// Workflows API
export const workflows = {
  create: async (data: {
    name: string;
    trigger_type: string;
    action_type: string;
    config_json: string;
  }) => {
    const response = await apiClient('/workflows', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.json();
  },
  
  list: async () => {
    const response = await apiClient('/workflows');
    return response.json();
  },
  
  toggle: async (id: string) => {
    const response = await apiClient(`/workflows/${id}/toggle`, {
      method: 'PUT',
    });
    return response.json();
  },
  
  delete: async (id: string) => {
    const response = await apiClient(`/workflows/${id}`, {
      method: 'DELETE',
    });
    return response;
  },
};

// Logs API
export const logs = {
  list: async (workflowId?: string) => {
    const endpoint = workflowId ? `/logs?workflow_id=${workflowId}` : '/logs';
    const response = await apiClient(endpoint);
    return response.json();
  },
};

