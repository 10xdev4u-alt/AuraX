import axios from 'axios'

const API_BASE_URL = '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const deviceService = {
  getAll: () => api.get('/devices'),
  getById: (id) => api.get(`/devices/${id}`),
  create: (data) => api.post('/devices', data),
}

export const firmwareService = {
  getAll: () => api.get('/firmware'),
  getById: (id) => api.get(`/firmware/${id}`),
  upload: (formData) => api.post('/firmware', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
}

export const releaseService = {
  getAll: () => api.get('/releases'),
  getById: (id) => api.get(`/releases/${id}`),
  create: (data) => api.post('/releases', data),
  updateStatus: (id, data) => api.put(`/releases/${id}/status`, data),
}

export default api
