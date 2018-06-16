import { handleResponse } from '@/requests/utils';

export function authenticate(identifier, password) {
  return fetch('/login', {
    method: "POST",
    body: JSON.stringify({id: identifier, password: password})
  })
    .then(handleResponse)
    .then(function (token) {
      return token.value;
    });
}
