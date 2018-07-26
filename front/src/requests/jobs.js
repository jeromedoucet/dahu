import { handleResponse } from '@/requests/utils';
import { getToken } from '@/requests/utils';

export function createJob(job) {
  return fetch('/jobs', {
    method: "POST",
    headers: { Authorization: getToken() },
    body: JSON.stringify(job)
  })
    .then(handleResponse);
}

export function fetchJobs() {
  return fetch('/jobs', {
    method: "GET",
    headers: { Authorization: getToken() }
  })
    .then(handleResponse);
}
