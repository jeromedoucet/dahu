import fetchMock from "fetch-mock";
import { 
  createDockerRegistry,
  fetchDockerRegistries,
  deleteDockerRegistry ,
  updateDockerRegistry,
} from '@/requests/dockerRegistries';
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';

describe('docker registries requests', () => {

  beforeEach(() => {
    fetchMock.restore();
    Cookies.remove(SESSION_COOKIE);
  });

  describe('create docker registries', () => {
    it('should send a post to /containers/docker/registries', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const name = 'some-name';
      const url = 'domaine/registry';
      const registry = {
        name,
        url
      };

      const dockerRegistryCreationMatcher = (callUrl, opt) => {
        return (
          callUrl === '/containers/docker/registries' &&
          opt.headers.Authorization === 'Bearer sometoken'  &&
          JSON.parse(opt.body).name === name &&
          JSON.parse(opt.body).url === url
        );
      };

      fetchMock.postOnce(dockerRegistryCreationMatcher, {
        status: 201
      });

      // when
      await expect(createDockerRegistry(registry)).resolves.toEqual({});
    });
  });

  describe('get docker registries', () => {
    it('should send a get to /containers/docker/registries', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const result = [{name: 'registry 1'}, {name: 'registry 2'}];

      const dockerRegistriesMatcher = (url, opt) => {
        return (
          url === '/containers/docker/registries' &&
          opt.headers.Authorization === 'Bearer sometoken' 
        );
      };

      fetchMock.getOnce(dockerRegistriesMatcher, {
        body: result
      });

      // when
      await expect(fetchDockerRegistries()).resolves.toEqual(result);
    });
  });

  describe('delete docker registry', () => {
    it('should send a delete to /containers/docker/registries/:id', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const id = "1";

      const dockerRegistryMatcher = (url, opt) => {
        return (
          url === `/containers/docker/registries/${id}` &&
          opt.headers.Authorization === 'Bearer sometoken' 
        );
      };

      fetchMock.deleteOnce(dockerRegistryMatcher, {
        status: 200
      });

      // when
      await expect(deleteDockerRegistry(id)).resolves.toEqual({});
    });
  });

  describe('update docker registry', () => {
    it('should send a delete to /containers/docker/registries/:id', async () => {
      // given
      Cookies.set(SESSION_COOKIE, 'sometoken');
      const id = '1';
      const name = 'some-name';
      const url = 'domaine/registry';
      const registry = {
        id,
        name,
        url,
      };

      const dockerRegistryMatcher = (callUrl, opt) => {
        return (
          callUrl === `/containers/docker/registries/${id}` &&
          opt.headers.Authorization === 'Bearer sometoken' &&
          JSON.parse(opt.body).name === name &&
          JSON.parse(opt.body).url === url
        );
      };

      fetchMock.putOnce(dockerRegistryMatcher, {
        body: registry
      });

      // when
      await expect(updateDockerRegistry(id, registry)).resolves.toEqual(registry);
    });
  });
});
