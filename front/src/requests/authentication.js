import { login } from '@/services/user';


function handleResponse(res) {
  if (!res.ok) {
    return Promise.reject(new Error(res.status));
  } else {
    return res.json();
  }
}

export function authenticate(identifier, password) {
  return fetch('/login', {
    method: "POST",
    body: JSON.stringify({id: identifier, password: password})
  })
    .then(handleResponse)
    .then(token => {
      login(token.value)
    })
  ;
}
