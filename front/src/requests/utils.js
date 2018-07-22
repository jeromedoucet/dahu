import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export class FetchError extends Error {
  constructor(msg, status) {
    super(msg);
    this.status = status;
  }
}

export async function handleResponse(res) {
  const body = await _parseJSON(res);
  if (!res.ok) {
    return Promise.reject(new FetchError(body.msg, res.status));
  } else {
    return body;
  }
}

function _parseJSON(response) {
  return response.text().then(function(text) {
    return text ? JSON.parse(text) : '';
  })
}

export function getToken() {
  return `Bearer ${Cookies.get(SESSION_COOKIE)}`;
}
