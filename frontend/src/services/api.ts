import type { Topic } from '../types/Topic';
import type { User, LoginCredentials, RegisterCredentials } from '../types/Auth';

const API_BASE_URL = 'http://localhost:8080/api';

export const fetchTopics = async (): Promise<Topic[]> => {
  const response = await fetch(`${API_BASE_URL}/topics`);
  
  const responseData = await response.json();
  
  if (!response.ok) {
    // If the response has an error message, use it; otherwise use a generic message
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }
  
  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
  
  return responseData.data || [];
};

// Authentication API functions
export const login = async (credentials: LoginCredentials): Promise<User> => {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  const responseData = await response.json();

  if (!response.ok) {
    // If the response has an error message, use it; otherwise use a generic message
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }

  // Store the token if provided
  if (responseData.token) {
    localStorage.setItem('auth_token', responseData.token);
  }

  return responseData.data;
};

export const register = async (credentials: RegisterCredentials): Promise<User> => {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  const responseData = await response.json();

  if (!response.ok) {
    // If the response has an error message, use it; otherwise use a generic message
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }

  // Store the token if provided
  if (responseData.token) {
    localStorage.setItem('auth_token', responseData.token);
  }

  return responseData.data;
};

export const logout = (): void => {
  localStorage.removeItem('auth_token');
};

export const getAuthToken = (): string | null => {
  return localStorage.getItem('auth_token');
};

export const isAuthenticated = (): boolean => {
  return !!getAuthToken();
};

// Voting API functions
export const vote = async (topicId: number, choice: 'Sim' | 'Não'): Promise<void> => {
  const token = getAuthToken();
  const response = await fetch(`${API_BASE_URL}/topics/${topicId}/vote`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({ choice }),
  });

  const responseData = await response.json();

  if (!response.ok) {
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
};

export const getVoteResult = async (topicId: number): Promise<{Sim: number, Não: number}> => {
  const response = await fetch(`${API_BASE_URL}/topics/${topicId}/result`);
  
  const responseData = await response.json();
  
  if (!response.ok) {
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }
  
  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
  
  return responseData.data;
};

export const openVotingSession = async (topicId: number, durationMinutes: number = 1): Promise<void> => {
  const token = getAuthToken();
  const response = await fetch(`${API_BASE_URL}/topics/${topicId}/session`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({ duration_minutes: durationMinutes }),
  });

  const responseData = await response.json();

  if (!response.ok) {
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
};
