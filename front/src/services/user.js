import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export function isAuthenticated() {
  return Cookies.get(SESSION_COOKIE) !== undefined;
}

export function logout() {
  Cookies.remove(SESSION_COOKIE);
}

export function login(token) {
  Cookies.set(SESSION_COOKIE, token);
}
