import fetchMock from "fetch-mock";
import { testRepo } from '@/requests/scm';
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('scm requests', () => {

  describe('test scm repo configuration', () => {

    it('should call the endpoint with auth token and correct body', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const url = 'ssh://git@localhost:10022/tester/test-repo.git';
      const key = 'some-private-key'
      const scmConf = {
        sshAuth: {
          url,
          key
        }
      };

      const scmMatcher = (callUrl, opt) => {
        return (
          callUrl === "/scm/git/repository" &&
          opt.headers.Authorization === 'Bearer sometoken'  &&
          JSON.parse(opt.body).sshAuth.url === url &&
          JSON.parse(opt.body).sshAuth.key === key
        );
      };

      fetchMock.postOnce(scmMatcher, {
        status: 200
      });

      // when
      await testRepo(scmConf);

      // then
      expect(fetchMock.called(scmMatcher, 'post')).toBeTruthy();
    });
  });
});
