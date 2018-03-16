import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export function isAuthenticated() {
  return Cookies.get(SESSION_COOKIE) !== undefined;
}
