import fetchMock from "fetch-mock";
import * as user from '@/services/user';
import { authenticate } from '@/requests/authentication';

user.login = jest.fn()

describe('authentication requests', () => {
  
  beforeEach(() => {
    fetchMock.restore();
    user.login.mockClear();
  });

  it('authenticate the user successfuly', async () => {
    const identifier = 'tester';
    const password = 'some-password'
    const postAttachmentMatcher = (url, opt) => {
      return (
        url === "/login" &&
        JSON.parse(opt.body).id === identifier &&
        JSON.parse(opt.body).password === password 
      );
    };
    fetchMock.postOnce(postAttachmentMatcher, {
      body: { value: "someToken" }
    });

    await authenticate(identifier, password);
    expect(user.login.mock.calls.length).toBe(1);
    expect(user.login.mock.calls[0][0]).toBe('someToken');
  });
  it('authenticate the user un-successfuly', async () => {
    const identifier = 'tester';
    const password = 'some-password'
    const postAttachmentMatcher = (url, opt) => {
      return (
        url === "/login" &&
        JSON.parse(opt.body).id === identifier &&
        JSON.parse(opt.body).password === password 
      );
    };
    fetchMock.postOnce(postAttachmentMatcher, 403);

    await expect(authenticate(identifier, password)).rejects.toThrow('403');
    expect(user.login.mock.calls.length).toBe(0);
  });
});
