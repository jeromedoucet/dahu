import fetchMock from "fetch-mock";
import { fetchJobs } from '@/requests/jobs';
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('jobs requests', () => {

  beforeEach(() => {
    fetchMock.restore();
    Cookies.remove(SESSION_COOKIE);
  });

  describe('get jobs', () => {
    it('should send a get to /jobs', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const result = [{name: 'job 1'}, {name: 'job 2'}];

      const jobsMatcher = (url, opt) => {
        return (
          url === "/jobs" &&
          opt.headers.Authorization === 'Bearer sometoken' 
        );
      };

      fetchMock.getOnce(jobsMatcher, {
        body: result
      });

      // when
      await expect(fetchJobs()).resolves.toEqual(result);;
    });
  });
});
