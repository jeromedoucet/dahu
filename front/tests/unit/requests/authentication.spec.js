import fetchMock from "fetch-mock";
import { authenticate } from '@/requests/authentication';
import { FetchError } from '@/requests/utils';

describe('authentication requests', () => {
  
  beforeEach(() => {
    fetchMock.restore();
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

    await expect(authenticate(identifier, password)).resolves.toEqual('someToken');
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

    await expect(authenticate(identifier, password)).rejects.toEqual(new FetchError('', 403));
  });
});
