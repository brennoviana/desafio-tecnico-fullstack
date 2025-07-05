import type { Topic } from '../types/Topic';

export const fetchTopics = async (): Promise<Topic[]> => {
  const response = await fetch('http://localhost:8080/topics');
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  const responseData = await response.json();
  
  if (responseData.status === 'error') {
    throw new Error(responseData.error);
  }
  
  return responseData.data || [];
};
