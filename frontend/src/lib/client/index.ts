/**
 * HTTP client for API communication
 */

export interface ApiResponse<T> {
data?: T;
error?: string;
message?: string;
}

export class ApiClient {
private baseUrl: string;
private token: string | null = null;

constructor(baseUrl = 'http://localhost:8080') {
this.baseUrl = baseUrl;
}

setAuthToken(token: string) {
this.token = token;
}

private getHeaders(): HeadersInit {
const headers: HeadersInit = {
'Content-Type': 'application/json'
};

if (this.token) {
headers.Authorization = `Bearer ${this.token}`;
}

return headers;
}

async get<T>(endpoint: string): Promise<T> {
const response = await fetch(`${this.baseUrl}${endpoint}`, {
method: 'GET',
headers: this.getHeaders()
});

if (!response.ok) {
throw new Error(`GET ${endpoint}: ${response.status} ${response.statusText}`);
}

return response.json();
}

async post<T>(endpoint: string, data?: any): Promise<T> {
const response = await fetch(`${this.baseUrl}${endpoint}`, {
method: 'POST',
headers: this.getHeaders(),
body: data ? JSON.stringify(data) : undefined
});

if (!response.ok) {
throw new Error(`POST ${endpoint}: ${response.status} ${response.statusText}`);
}

return response.json();
}

async put<T>(endpoint: string, data: any): Promise<T> {
const response = await fetch(`${this.baseUrl}${endpoint}`, {
method: 'PUT',
headers: this.getHeaders(),
body: JSON.stringify(data)
});

if (!response.ok) {
throw new Error(`PUT ${endpoint}: ${response.status} ${response.statusText}`);
}

return response.json();
}

async delete(endpoint: string): Promise<void> {
const response = await fetch(`${this.baseUrl}${endpoint}`, {
method: 'DELETE',
headers: this.getHeaders()
});

if (!response.ok) {
throw new Error(`DELETE ${endpoint}: ${response.status} ${response.statusText}`);
}
}
}

// Global instance
export const apiClient = new ApiClient();
