import { checkAuthenticationBeforeNavigation } from '@/routes'
import * as user from '@/services/user';

describe('routes', () => {

  describe('checkAuthenticationBeforeNavigation', () => {

    user.isAuthenticated = jest.fn();
    const next = jest.fn();

    beforeEach(() => {
      user.isAuthenticated.mockClear();
      next.mockClear();
    });

    it('should redirect to login if the destination is not Login and if the user is note authenticated', () => {
      // given
      const to = { path: '/'};
      const from = {};
      user.isAuthenticated.mockReturnValue(false);

      // when
      checkAuthenticationBeforeNavigation(to, from, next);

      // then
      expect(next.mock.calls.length).toBe(1);
      expect(next.mock.calls[0][0]).toBe('/login');
    });

    it('should allow navigation to Home if the user is authenticated', () => {
      // given
      const to = { path: '/'};
      const from = {};
      user.isAuthenticated.mockReturnValue(true);

      // when
      checkAuthenticationBeforeNavigation(to, from, next);

      // then
      expect(next.mock.calls.length).toBe(1);
      expect(next.mock.calls[0].length).toBe(0);
    });

    it('should allow navigation if the user is not authenticated but the destination is Login', () => {
      // given
      const to = { path: '/login'};
      const from = {};
      user.isAuthenticated.mockReturnValue(false);

      // when
      checkAuthenticationBeforeNavigation(to, from, next);

      // then
      expect(next.mock.calls.length).toBe(1);
      expect(next.mock.calls[0].length).toBe(0);
    });

    it('should allow redirect to home if the user is authenticated but the destination is Login', () => {
      // given
      const to = { path: '/login'};
      const from = {};
      user.isAuthenticated.mockReturnValue(true);

      // when
      checkAuthenticationBeforeNavigation(to, from, next);

      // then
      expect(next.mock.calls.length).toBe(1);
      expect(next.mock.calls[0][0]).toBe('/');
    });
  })
});
