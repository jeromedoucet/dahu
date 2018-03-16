import fetchMock from "fetch-mock";
import * as user from '@/services/user';
import { authenticate } from '@/requests/authentication';

user.login = jest.fn()

describe('authentication requests', () => {

  it('authenticate the user sucessfuly', async () => {
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
});
