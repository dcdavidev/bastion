import axios from 'axios';
import Cookies from 'js-cookie';

export const api = axios.create({
  baseURL: '/api/v1',
});

// Add interceptor to include Bearer token from cookie if available
api.interceptors.request.use((config) => {
  const token = Cookies.get('bastion_session');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const getApi = (request?: Request) => {
  if (request) {
    const cookieHeader = request.headers.get('Cookie') || '';
    const token = cookieHeader.match(/bastion_session=([^;]+)/)?.[1];
    
    // Server-side axios needs absolute URL
    return axios.create({
      baseURL: 'http://localhost:8287/api/v1',
      headers: {
        Authorization: token ? `Bearer ${token}` : '',
      },
    });
  }
  return api;
};
