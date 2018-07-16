import { handleResponse } from '@/requests/utils';
import { getToken } from '@/requests/utils';

export function fetchJobs() {
  return fetch('/jobs', {
    method: "GET",
    headers: { Authorization: getToken() }
  })
    .then(handleResponse);
}
