import { handleResponse } from '@/requests/utils';
import { getToken } from '@/requests/utils';

export function createDockerRegistry(registry) {
  return fetch('/containers/docker/registries', {
    method: "POST",
    headers: { Authorization: getToken() },
    body: JSON.stringify(registry)
  })
    .then(handleResponse);
}

export function fetchDockerRegistries() {
  return fetch('/containers/docker/registries', {
    method: "GET",
    headers: { Authorization: getToken() }
  })
    .then(handleResponse);
}

export function deleteDockerRegistry(id) {
  return fetch(`/containers/docker/registries/${id}`, {
    method: "DELETE",
    headers: { Authorization: getToken() }
  })
    .then(handleResponse);
}

export function updateDockerRegistry(id, registry) {
  return fetch(`/containers/docker/registries/${id}`, {
    method: "PUT",
    headers: { Authorization: getToken() },
    body: JSON.stringify(registry)
  })
    .then(handleResponse);
}

export function testDockerRegistry(registry) {
  return fetch('/containers/docker/registries/test', {
    method: "POST",
    headers: { Authorization: getToken() },
    body: JSON.stringify(registry)
  })
    .then(handleResponse);
}
