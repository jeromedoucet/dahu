import { mount } from '@vue/test-utils';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import VueRouter from 'vue-router';
import NavBar from '@/components/NavBar.vue'
import { SESSION_COOKIE } from '@/constants'
import Cookies from 'js-cookie';
import { isAuthenticated } from '@/services/user';

// tested component conf
const localVue = createLocalVue();
localVue.use(BootstrapVue);
localVue.use(VueRouter);

const router = new VueRouter({
  routes: []
})

const createNavBar = propsData => mount(NavBar, { propsData, localVue, router });

describe('NavBar.vue', () => {

  beforeEach(() => {
    Cookies.remove(SESSION_COOKIE);
    jest.clearAllMocks();
  });

  describe('disconnect', () => {

    it("should logout and redirect to '/login' when the user is login", () => {
      // given
      Cookies.set(SESSION_COOKIE, 'token');
      let navBarCmp = createNavBar();
      navBarCmp.vm.$router.push = jest.fn();

      // when
      navBarCmp.vm.disconnect();

      // then
      expect(isAuthenticated()).toBe(false);
      expect(navBarCmp.vm.$router.push.mock.calls.length).toBe(1);
      expect(navBarCmp.vm.$router.push.mock.calls[0][0]).toBe('/login');
    });

  })
});
