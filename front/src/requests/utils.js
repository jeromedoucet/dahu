import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export function handleResponse(res) {
  if (!res.ok) {
    return Promise.reject(new Error(res.status));
  } else {
    return res.json();
  }
}

export function getToken() {
  return `Bearer ${Cookies.get(SESSION_COOKIE)}`;
}
