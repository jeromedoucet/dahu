import { 
  isAuthenticated,
  logout,
  login
} from '@/services/user'
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('user services', () => {

  beforeEach(() => {
    Cookies.remove(SESSION_COOKIE);
  });

  describe('isAuthenticated', () => {

    it('return true if a session cookie is present', () => {
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const authenticated = isAuthenticated();
      expect(authenticated).toBeTruthy();
    });

    it('return false if a session cookie is not present', () => {
      const authenticated = isAuthenticated();
      expect(authenticated).toBeFalsy();
    });
  });

  describe('logout', () => {

    it('remove session cookie', () => {
      Cookies.set(SESSION_COOKIE, 'sometoken');
      logout();
      expect(Cookies.get(SESSION_COOKIE)).toBeUndefined();
    });
  });

  describe('login', () => {
    
    it('set session cookie', () => {
      login('someToken');
      expect(Cookies.get(SESSION_COOKIE)).toBe('someToken');
    })
  });
});
