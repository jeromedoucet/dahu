import { isAuthenticated } from '@/services/user'
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('isAuthenticated', () => {

  beforeEach(() => {
    Cookies.remove(SESSION_COOKIE);
  });

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
