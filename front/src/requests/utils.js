import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export class FetchError {
  constructor(msg, status) {
    this.status = status;
    this.message = msg;
  }
}

export async function handleResponse(res) {
  const body = await _parseJSON(res);
  if (!res.ok) {
    return Promise.reject(new FetchError(body, res.status));
  } else {
    return body;
  }
}

function _parseJSON(response) {
  return response.text().then(function(text) {
    return text ? JSON.parse(text) : {};
  })
}

export function getToken() {
  return `Bearer ${Cookies.get(SESSION_COOKIE)}`;
}
