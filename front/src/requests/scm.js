import { handleResponse } from '@/requests/utils';
import { getToken } from '@/requests/utils';

export function testRepo(repoConfig) {
  return fetch('/scm/git/repository', {
    method: "POST",
    headers: { Authorization: getToken() },
    body: JSON.stringify(repoConfig)
  })
    .then(handleResponse);
}
