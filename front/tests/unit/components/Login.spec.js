import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import fetchMock from "fetch-mock";
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import flushPromises from 'flush-promises';
import { isAuthenticated } from '@/services/user';
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';
import Login from '@/components/Login.vue'
 
// tested component conf
const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
})

const createLogin = propsData => mount(Login, { propsData, localVue, router });

let loginCmp = null;

// fetch mock configuration
const identifier = 'tester';
const password = 'test';
const badIdentifier = 'badTester';
const badPassword = 'badTest';
const token = 'someToken';

const loginSuccessMatcher = (url, opt) => {
  return (
    url === "/login" &&
    JSON.parse(opt.body).id === identifier &&
    JSON.parse(opt.body).password === password
  );
};

const login401Matcher = (url, opt) => {
  return (
    url === "/login" &&
    JSON.parse(opt.body).id === identifier &&
    JSON.parse(opt.body).password === badPassword
  );
};

const login404Matcher = (url, opt) => {
  return (
    url === "/login" &&
    JSON.parse(opt.body).id === badIdentifier
  );
};

const login500Matcher = (url, opt) => {
  return (
    url === "/login" &&
    !JSON.parse(opt.body).id &&
    !JSON.parse(opt.body).password
  );
};

fetchMock.postOnce(loginSuccessMatcher, {
  body: { value: token }
});

fetchMock.postOnce(login401Matcher, {
  status: 401
});

fetchMock.postOnce(login404Matcher, {
  status: 404
});

fetchMock.postOnce(login500Matcher, {
  status: 500
});

describe('Login.vue', () => {

  afterEach(() => {
    jest.clearAllMocks();
    Cookies.remove(SESSION_COOKIE);
  });

  it('should emit on-login on click when identifier and password are set', async () => {
    // given
    loginCmp = createLogin();
    loginCmp.setData({ form:{ identifier: identifier, password: password }});
    loginCmp.vm.$router.go = jest.fn();
    const btn = loginCmp.find('#login-submit-button');
    const evt = { preventDefault: jest.fn() };   

    // when
    loginCmp.vm.onSubmit(evt);
    await flushPromises();

    // then
    expect(isAuthenticated()).toBe(true);
    expect(loginCmp.vm.$router.go.mock.calls.length).toBe(1);
    expect(loginCmp.vm.$router.go.mock.calls[0][0]).toBe('/');
  });

  it('should print the right error message when unknown user', async () => {
    // given
    loginCmp = createLogin();
    loginCmp.setData({ form:{ identifier: badIdentifier, password: badPassword }});
    loginCmp.vm.$router.go = jest.fn();
    const btn = loginCmp.find('#login-submit-button');
    const evt = { preventDefault: jest.fn() };   

    // when
    loginCmp.vm.onSubmit(evt);
    await flushPromises();

    // then
    expect(isAuthenticated()).toBe(false);
    expect(loginCmp.vm.$router.go.mock.calls.length).toBe(0);
    expect(loginCmp.vm.errorMessage).toBe('Authentication error. Please check your credentials and try again.');
  });

  it('should print the right error message when bad password', async () => {
    // given
    loginCmp = createLogin();
    loginCmp.setData({ form:{ identifier: identifier, password: badPassword }});
    loginCmp.vm.$router.go = jest.fn();
    const btn = loginCmp.find('#login-submit-button');
    const evt = { preventDefault: jest.fn() };   

    // when
    loginCmp.vm.onSubmit(evt);
    await flushPromises();

    // then
    expect(isAuthenticated()).toBe(false);
    expect(loginCmp.vm.$router.go.mock.calls.length).toBe(0);
    expect(loginCmp.vm.errorMessage).toBe('Authentication error. Please check your credentials and try again.');
  });

  it('should print the right error message when unknow error', async () => {
    // given
    loginCmp = createLogin();
    loginCmp.setData({ form:{ identifier: null, password: null }});
    loginCmp.vm.$router.go = jest.fn();
    const btn = loginCmp.find('#login-submit-button');
    const evt = { preventDefault: jest.fn() };   

    // when
    loginCmp.vm.onSubmit(evt);
    await flushPromises();

    // then
    expect(isAuthenticated()).toBe(false);
    expect(loginCmp.vm.$router.go.mock.calls.length).toBe(0);
    expect(loginCmp.vm.errorMessage).toBe('Unknown error. Please contact your administrator.');
  });

  // TODO : forgotten password ?
})
