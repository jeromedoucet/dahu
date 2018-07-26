import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import * as scm from '@/requests/scm';
import * as jobsRequest from '@/requests/jobs';
import flushPromises from 'flush-promises'
import { FetchError } from '@/requests/utils'
import JobCreation from '@/components/jobs/JobCreation.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
});

const createJobCreation = propsData => mount(JobCreation, { propsData, localVue, router });

describe('JobCreation', () => {

  beforeAll(() => {
    scm.testRepo = jest.fn();
    jobsRequest.createJob = jest.fn();
  });

  beforeEach(() => {
    scm.testRepo.mockReset(); 
    jobsRequest.createJob.mockReset();
  });

  afterAll(() => {
    scm.testRepo.mockRestore(); 
    jobsRequest.createJob.mockRestore();
  });

  describe('isFormValid', () => {
    it('should return false when no auth scheme', () => {
      // given
      const cmp = createJobCreation();
      cmp.setData({name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeFalsy();
    });

    it('should return false when no name', () => {
      // given
      const scheme = 'http';
      const httpForm = {
        url: 'https://github.com:test/test-repo.git',
        user: 'some-user',
        password: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, httpForm});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeFalsy();
    });

    it('should return false when missing url on http scheme', () => {
      // given
      const scheme = 'http';
      const httpForm = {
        user: 'some-user',
        password: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, httpForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeFalsy();
    });

    it('should return true when url only on http scheme', () => {

      // given
      const scheme = 'http';
      const httpForm = {
        url: 'https://github.com:test/test-repo.git',
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, httpForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeTruthy();
    });

    it('should return true when full http scheme', () => {

      // given
      const scheme = 'http';
      const httpForm = {
        url: 'https://github.com:test/test-repo.git',
        user: 'some-user',
        password: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, httpForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeTruthy();
    });

    it('should return true when full ssh scheme', () => {
      // given
      const scheme = 'ssh';
      const sshForm = {
        url: 'git@github.com:test/test-repo.git',
        key: 'some-private-key',
        keyPassword: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, sshForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeTruthy();
    });

    it('should return true when minimal ssh scheme', () => {
      // given
      const scheme = 'ssh';
      const sshForm = {
        url: 'git@github.com:test/test-repo.git',
        key: 'some-private-key',
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, sshForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeTruthy();
    });

    it('should return false when missing url on ssh scheme', () => {
      // given
      const scheme = 'ssh';
      const sshForm = {
        key: 'some-private-key',
        keyPassword: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, sshForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeFalsy();
    });

    it('should return false when missing key on ssh scheme', () => {
      // given
      const scheme = 'ssh';
      const sshForm = {
        url: 'git@github.com:test/test-repo.git',
        keyPassword: 'some-password'
      };
      const cmp = createJobCreation();
      cmp.setData({authSchemSelected: scheme, sshForm, name: "some job"});

      // when
      const isFormValid = cmp.vm.isFormValid;

      // then
      expect(isFormValid).toBeFalsy();
    });

  });

  it('should has no auth scheme configuration selected by default', () => {
    // when
    const cmp = createJobCreation();

    // then
    expect(cmp.find('#http-auth-conf').exists()).toBeFalsy();
    expect(cmp.find('#ssh-auth-conf').exists()).toBeFalsy();
  });

  it('should switch to http configuration scheme', () => {
    // given
    const scheme = 'http';

    // when
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme});

    // then
    expect(cmp.find('#http-auth-conf').exists()).toBeTruthy();
    expect(cmp.find('#ssh-auth-conf').exists()).toBeFalsy();
  });

  it('should switch to ssh configuration scheme', () => {
    // given
    const scheme = 'ssh';

    // when
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme});

    // then
    expect(cmp.find('#http-auth-conf').exists()).toBeFalsy();
    expect(cmp.find('#ssh-auth-conf').exists()).toBeTruthy();
  });

  it('should display alert success when auth conf test is successful on ssh', async () => {
    // given
    const scheme = 'ssh';
    const sshForm = {
      url: 'git@github.com:test/test-repo.git',
      key: 'some-private-key',
      keyPassword: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, sshForm});
    scm.testRepo.mockResolvedValue(200);

    // when
    cmp.find('#job-test-button').trigger('click');
    await flushPromises();

    // then
    expect(scm.testRepo.mock.calls.length).toBe(1);
    expect(scm.testRepo.mock.calls[0][0]).toEqual({sshAuth: sshForm });
    expect(cmp.find('#repo-test-success-msg').exists()).toBeTruthy();
  });

  it('should display alert failure when auth conf test has failed on ssh', async () => {
    // given
    const scheme = 'ssh';
    const sshForm = {
      url: 'git@github.com:test/test-repo.git',
      key: 'some-private-key',
      keyPassword: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, sshForm});
    scm.testRepo.mockRejectedValue(new FetchError('some error', 400));

    // when
    cmp.find('#job-test-button').trigger('click');
    await flushPromises();

    // then
    expect(scm.testRepo.mock.calls.length).toBe(1);
    expect(scm.testRepo.mock.calls[0][0]).toEqual({sshAuth: sshForm });
    expect(cmp.find('#repo-test-failure-msg').exists()).toBeTruthy();
    expect(cmp.find('#repo-test-failure-msg').text()).toBe('An error has happend during test : some error');
  });

  it('should display alert success when auth conf test is successful on http', async () => {
    // given
    const scheme = 'http';
    const httpForm = {
      url: 'https://github.com:test/test-repo.git',
      user: 'some-user',
      password: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, httpForm});
    scm.testRepo.mockResolvedValue(200);

    // when
    cmp.find('#job-test-button').trigger('click');
    await flushPromises();

    // then
    expect(scm.testRepo.mock.calls.length).toBe(1);
    expect(scm.testRepo.mock.calls[0][0]).toEqual({httpAuth: httpForm });
    expect(cmp.find('#repo-test-success-msg').exists()).toBeTruthy();
  });

  it('should display alert failure when auth conf test has failed on http', async () => {
    // given
    const scheme = 'http';
    const httpForm = {
      url: 'https://github.com:test/test-repo.git',
      user: 'some-user',
      password: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, httpForm});
    scm.testRepo.mockRejectedValue(new FetchError('some error', 400));

    // when
    cmp.find('#job-test-button').trigger('click');
    await flushPromises();

    // then
    expect(scm.testRepo.mock.calls.length).toBe(1);
    expect(scm.testRepo.mock.calls[0][0]).toEqual({httpAuth: httpForm });
    expect(cmp.find('#repo-test-failure-msg').exists()).toBeTruthy();
    expect(cmp.find('#repo-test-failure-msg').text()).toBe('An error has happend during test : some error');
  });

  it('should redirect to / when creation is successful', async () => {
    // given
    const scheme = 'http';
    const httpForm = {
      url: 'https://github.com:test/test-repo.git',
      user: 'some-user',
      password: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, httpForm, name: 'some name'});
    jobsRequest.createJob.mockResolvedValue(201);
    cmp.vm.$router.go = jest.fn();

    // when
    cmp.find('#job-creation-button').trigger('click');
    await flushPromises();

    // then
    expect(jobsRequest.createJob.mock.calls.length).toBe(1);
    expect(jobsRequest.createJob.mock.calls[0][0]).toEqual({name: 'some name', gitConfig: {httpAuth: httpForm }});
    expect(cmp.vm.$router.go.mock.calls.length).toBe(1);
    expect(cmp.vm.$router.go.mock.calls[0][0]).toBe('/');
  });

  it('should display error message when creation has failed', async () => {
    // given
    const scheme = 'http';
    const httpForm = {
      url: 'https://github.com:test/test-repo.git',
      user: 'some-user',
      password: 'some-password'
    };
    const cmp = createJobCreation();
    cmp.setData({authSchemSelected: scheme, httpForm, name: 'some name'});
    jobsRequest.createJob.mockRejectedValue(new FetchError('some error', 400));
    cmp.vm.$router.go = jest.fn();

    // when
    cmp.find('#job-creation-button').trigger('click');
    await flushPromises();

    // then
    expect(jobsRequest.createJob.mock.calls.length).toBe(1);
    expect(jobsRequest.createJob.mock.calls[0][0]).toEqual({name: 'some name', gitConfig: {httpAuth: httpForm }});
    expect(cmp.find('#repo-test-failure-msg').exists()).toBeTruthy();
    expect(cmp.find('#repo-test-failure-msg').text()).toBe('An error has happend during the creation : some error');
  });

});
