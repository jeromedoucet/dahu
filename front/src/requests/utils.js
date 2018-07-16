import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export function handleResponse(res) {
  if (!res.ok) {
    return Promise.reject(new Error(res.status));
  } else {
    return _parseJSON(res);
  }
}

function _parseJSON(response) {
  return response.text().then(function(text) {
    return text ? JSON.parse(text) : {}
  })
}

export function getToken() {
  return `Bearer ${Cookies.get(SESSION_COOKIE)}`;
}
