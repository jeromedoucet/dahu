import { mount, createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import * as registriesRequests from '@/requests/dockerRegistries';
import flushPromises from 'flush-promises'
import { FetchError } from '@/requests/utils'
import RegistryDetails from '@/components/registries/RegistryDetails.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
});

const createRegistryDetails = propsData => mount(RegistryDetails, { propsData, localVue, router });

describe('RegistryDetails', () => {

  beforeAll(() => {
    registriesRequests.testDockerRegistry = jest.fn();
    registriesRequests.createDockerRegistry = jest.fn();
		registriesRequests.updateDockerRegistry = jest.fn();
  });

  beforeEach(() => {
    registriesRequests.testDockerRegistry.mockReset();
    registriesRequests.createDockerRegistry.mockReset();
    registriesRequests.updateDockerRegistry.mockReset();
  });

  afterAll(() => {
    registriesRequests.testDockerRegistry.mockRestore();
    registriesRequests.createDockerRegistry.mockRestore();
    registriesRequests.updateDockerRegistry.mockRestore();
  });

  describe('test configuration', () => {
    it('should print success message when 200', async () => {
      // given
      const registry = {
        name: 'registry for test',
        url: 'https://hub.docker.com',
        password: '',
        user: '',
      };
      const cmp = createRegistryDetails({registry});
      registriesRequests.testDockerRegistry.mockResolvedValue(200);

      // when
      cmp.find('#registry-test-button').trigger('click');
      await flushPromises();

      // then
      expect(registriesRequests.testDockerRegistry.mock.calls.length).toBe(1);
      expect(registriesRequests.testDockerRegistry.mock.calls[0][0]).toEqual(registry);
      expect(cmp.find('#registry-success-msg').exists()).toBeTruthy();
      expect(cmp.find('#registry-success-msg').text()).toBe('The test is successful');
    });

    it('should print success message when 200', async () => {
      // given
      const registry = {
        name: 'registry for test',
        url: 'https://hub.docker.com',
        password: '',
        user: '',
      };
      const cmp = createRegistryDetails({registry});
      registriesRequests.testDockerRegistry.mockRejectedValue(new FetchError({ msg: 'some error'}, 400));

      // when
      cmp.find('#registry-test-button').trigger('click');
      await flushPromises();

      // then
      expect(registriesRequests.testDockerRegistry.mock.calls.length).toBe(1);
      expect(registriesRequests.testDockerRegistry.mock.calls[0][0]).toEqual(registry);
      expect(cmp.find('#registry-failure-msg').exists()).toBeTruthy();
      expect(cmp.find('#registry-failure-msg').text()).toBe('An error has happend during test : some error');
    });
  });
  describe('creation', () => {
    it('should send a creation request when at least url and name are non null and registry props is null', async () => {
       // given
      const registry = {
        name: 'registry for test',
        url: 'https://hub.docker.com',
        password: '',
        user: '',
      };
      const stub = jest.fn()
      const cmp = createRegistryDetails();
      cmp.setData({form: registry});
			cmp.vm.$on('registry-saved', stub)
      registriesRequests.createDockerRegistry.mockResolvedValue(200);

      // when
      cmp.find('#registry-save-button').trigger('click');
      await flushPromises();   

      expect(registriesRequests.createDockerRegistry.mock.calls.length).toBe(1);
      expect(registriesRequests.createDockerRegistry.mock.calls[0][0]).toEqual(registry);
			expect(stub).toHaveBeenCalled();
    });

    it('should disable creation button when no name', async () => {
       // given
      const registry = {
        name: '',
        url: 'https://hub.docker.com',
        password: '',
        user: '',
      };
      const stub = jest.fn()
      const cmp = createRegistryDetails();
      cmp.setData({form: registry});
			cmp.vm.$on('registry-saved', stub)
      registriesRequests.createDockerRegistry.mockResolvedValue(200);

      // when
      cmp.find('#registry-save-button').trigger('click');
      await flushPromises();   

      expect(registriesRequests.createDockerRegistry.mock.calls.length).toBe(0);
			expect(stub).toHaveBeenCalledTimes(0);
    });

    it('should disable creation button when no url', async () => {
       // given
      const registry = {
        name: 'registry test',
        url: '',
        password: '',
        user: '',
      };
      const stub = jest.fn()
      const cmp = createRegistryDetails();
      cmp.setData({form: registry});
			cmp.vm.$on('registry-saved', stub)
      registriesRequests.createDockerRegistry.mockResolvedValue(200);

      // when
      cmp.find('#registry-save-button').trigger('click');
      await flushPromises();   

      expect(registriesRequests.createDockerRegistry.mock.calls.length).toBe(0);
			expect(stub).toHaveBeenCalledTimes(0);
    });

    it('should display error msg when error on registry creation', async () => {
       // given
      const registry = {
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: '',
      };
      const stub = jest.fn()
      const cmp = createRegistryDetails();
      cmp.setData({form: registry});
			cmp.vm.$on('registry-saved', stub)
      registriesRequests.createDockerRegistry.mockRejectedValue(new FetchError({msg: 'some error'}, 400));

      // when
      cmp.find('#registry-save-button').trigger('click');
      await flushPromises();   

      expect(registriesRequests.createDockerRegistry.mock.calls.length).toBe(1);
      expect(registriesRequests.createDockerRegistry.mock.calls[0][0]).toEqual(registry);
			expect(stub).toHaveBeenCalledTimes(0);
      expect(cmp.find('#registry-failure-msg').exists()).toBeTruthy();
      expect(cmp.find('#registry-failure-msg').text()).toBe('An error has happend during creation : some error');
    });
  });

	describe('update', () => {
    it('should update existing registry when it has changed', async () => {
       // given
      const initialRegistry = {
				id: 1,
				lastModificationTime: 123465,
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
      const registryForm = {
				lastModificationTime: 123465,
        name: 'registry test updated',
        url: 'https://hub.docker.com',
        password: 'test',
        user: 'tester',
      };
			const registryUpdate = {
				...registryForm,
				lastModificationTime: initialRegistry.lastModificationTime,
				changedFields: ["name", "password"]
			}
			const stub = jest.fn()
			const cmp = createRegistryDetails({registry: initialRegistry});
			cmp.setData({form: registryForm});
			cmp.vm.$on('registry-saved', stub)
			registriesRequests.updateDockerRegistry.mockResolvedValue(200);

			// when
			cmp.find('#registry-save-button').trigger('click');
			await flushPromises();   

			expect(registriesRequests.updateDockerRegistry.mock.calls.length).toBe(1);
			expect(registriesRequests.updateDockerRegistry.mock.calls[0][0]).toEqual(1);
			expect(registriesRequests.updateDockerRegistry.mock.calls[0][1]).toEqual(registryUpdate);
			expect(stub).toHaveBeenCalled();
    });

    it('should disable save button when nothing changed', async () => {
       // given
      const initialRegistry = {
				id: 1,
				lastModificationTime: 123465,
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
      const registryForm = {
				lastModificationTime: 123465,
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
			const stub = jest.fn()
			const cmp = createRegistryDetails({registry: initialRegistry});
			cmp.setData({form: registryForm});
			cmp.vm.$on('registry-saved', stub)

			// when
			cmp.find('#registry-save-button').trigger('click');
			await flushPromises();   

			expect(registriesRequests.updateDockerRegistry.mock.calls.length).toBe(0);
			expect(stub).toHaveBeenCalledTimes(0);
    });

    it('should replace form data when conflict', async () => {
       // given
      const initialRegistry = {
				id: 1,
				lastModificationTime: 123465,
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
      const registryForm = {
				lastModificationTime: 123465,
        name: 'registry test 2',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
			const updatedRegistry = {
				lastModificationTime: 123466,
        name: 'registry test 1',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
			}
			const stub = jest.fn()
			const cmp = createRegistryDetails({registry: initialRegistry});
			cmp.setData({form: registryForm});
			cmp.vm.$on('registry-saved', stub)
			registriesRequests.updateDockerRegistry.mockRejectedValue(new FetchError(updatedRegistry, 409));

			// when
			cmp.find('#registry-save-button').trigger('click');
			await flushPromises();   

			expect(registriesRequests.updateDockerRegistry.mock.calls.length).toBe(1);
			expect(stub).toHaveBeenCalledTimes(0);
      expect(cmp.find('#registry-failure-msg').exists()).toBeTruthy();
      expect(cmp.find('#registry-failure-msg').text()).toBe('An error has happend during the saving : there is a conflict ! the registry has been reloaded');
			expect(cmp.vm.$data.form).toBe(updatedRegistry);
    });

    it('should print an error message when issue append', async () => {
       // given
      const initialRegistry = {
				id: 1,
				lastModificationTime: 123465,
        name: 'registry test',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
      const registryForm = {
				lastModificationTime: 123465,
        name: 'registry test 2',
        url: 'https://hub.docker.com',
        password: '',
        user: 'tester',
      };
			const stub = jest.fn()
			const cmp = createRegistryDetails({registry: initialRegistry});
			cmp.setData({form: registryForm});
			cmp.vm.$on('registry-saved', stub)
			registriesRequests.updateDockerRegistry.mockRejectedValue(new FetchError({ msg: 'some error' }, 400));

			// when
			cmp.find('#registry-save-button').trigger('click');
			await flushPromises();   

			expect(registriesRequests.updateDockerRegistry.mock.calls.length).toBe(1);
			expect(stub).toHaveBeenCalledTimes(0);
      expect(cmp.find('#registry-failure-msg').exists()).toBeTruthy();
      expect(cmp.find('#registry-failure-msg').text()).toBe('An error has happend during the saving : some error');
    });
	});
});
