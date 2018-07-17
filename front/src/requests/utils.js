import Cookies from 'js-cookie';
import { SESSION_COOKIE } from '@/constants'

export async function handleResponse(res) {
  const body = await _parseJSON(res);
  if (!res.ok) {
    return Promise.reject(new Error(body.msg));
  } else {
    return body;
  }
}

function _parseJSON(response) {
  return response.text().then(function(text) {
    return text ? JSON.parse(text) : response.status
  })
}

export function getToken() {
  return `Bearer ${Cookies.get(SESSION_COOKIE)}`;
}
