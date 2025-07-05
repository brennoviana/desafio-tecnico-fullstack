import type { Topic } from '../types/Topic';
import type { User, LoginCredentials, RegisterCredentials } from '../types/Auth';

const API_BASE_URL = 'http://localhost:8080/api';

export const fetchTopics = async (): Promise<Topic[]> => {
  const response = await fetch(`${API_BASE_URL}/topics`);
  
  const responseData = await response.json();
  
  if (!response.ok) {
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }
  
  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
  
  return responseData.data || [];
};

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
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }

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
    throw new Error(responseData.error || `HTTP error! status: ${response.status}`);
  }

  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }

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

export const getSessionStatus = async (topicId: number): Promise<{open_at: number, close_at: number} | null> => {
  try {
    const response = await fetch(`${API_BASE_URL}/topics/${topicId}/session`);
    
    if (!response.ok) {
      return null;
    }
    
    const responseData = await response.json();
    
    if (responseData.status === 'error') {
      return null;
    }
    
    return responseData.data;
  } catch {
    return null;
  }
};
