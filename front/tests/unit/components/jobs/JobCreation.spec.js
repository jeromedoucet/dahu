import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import * as scm from '@/requests/scm';
import flushPromises from 'flush-promises'
import { FetchError } from '@/requests/utils'
import JobCreation from '@/components/jobs/JobCreation.vue'

const localVue = createLocalVue();
localVue.use(BootstrapVue);

const createJobCreation = propsData => mount(JobCreation, { propsData, localVue });

describe('JobCreation', () => {

  beforeAll(() => {
    scm.testRepo = jest.fn();
  });

  beforeEach(() => {
    scm.testRepo.mockReset(); 
  });

  afterAll(() => {
    scm.testRepo.mockRestore(); 
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
});
// todo test job creation
