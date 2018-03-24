import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import fetchMock from "fetch-mock";
import BootstrapVue from 'bootstrap-vue';
import flushPromises from 'flush-promises';
import { isAuthenticated } from '@/services/user';
import Login from '@/components/Login.vue'
 

const localVue = createLocalVue()
localVue.use(BootstrapVue)

const createLogin = propsData => mount(Login, { propsData, localVue });

let loginCmp = null;

const identifier = 'tester';
const password = 'test';
const token = 'someToken';

const postAttachmentMatcher = (url, opt) => {
  return (
    url === "/login" &&
    JSON.parse(opt.body).id === identifier &&
    JSON.parse(opt.body).password === password 
  );
};

fetchMock.postOnce(postAttachmentMatcher, {
  body: { value: token }
});

describe('Login.vue', () => {

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should emit on-login on click when identifier and password are set', async () => {
    // given
    loginCmp = createLogin();
    loginCmp.setData({ form:{ identifier: identifier, password: password }});
    const btn = loginCmp.find('#login-submit-button');
    const evt = { preventDefault: jest.fn() };   

    // when
    loginCmp.vm.onSubmit(evt);
    await flushPromises();

    // then
    expect(isAuthenticated()).toBe(true);
  });

  // TODO : validation
  // TODO : forgotten password ?
})
