import fetchMock from "fetch-mock";
import { fetchJobs, createJob } from '@/requests/jobs';
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('jobs requests', () => {

  beforeEach(() => {
    fetchMock.restore();
    Cookies.remove(SESSION_COOKIE);
  });

  describe('create jobs', () => {
    it('should send a post to /jobs', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const name = 'some-name';
      const url = "git@domaine/repo";
      const key = "some-private-key";
      const job = {
        name,
        gitConfig: {
          sshAuth: {
            url,
            key
          }
        }
      };

      const jobCreationMatcher = (callUrl, opt) => {
        return (
          callUrl === "/jobs" &&
          opt.headers.Authorization === 'Bearer sometoken'  &&
          JSON.parse(opt.body).name === name &&
          JSON.parse(opt.body).gitConfig.sshAuth.url === url &&
          JSON.parse(opt.body).gitConfig.sshAuth.key === key
        );
      };

      fetchMock.postOnce(jobCreationMatcher, {
        status: 201
      });

      // when
      await expect(createJob(job)).resolves.toEqual({});
    });
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
      await expect(fetchJobs()).resolves.toEqual(result);
    });
  });
});
