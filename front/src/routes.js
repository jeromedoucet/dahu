import Home from '@/components/Home.vue';
import Login from '@/components/Login.vue';

import { isAuthenticated } from '@/services/user';

export const routes = [
  { path: '/login', component: Login },
  { path: '/', component: Home }
]

export function checkAuthenticationBeforeNavigation(to, from, next) {
  if(isAuthenticated() || to.path === '/login') {
    next();
  } else {
    next('/login');
  }
}
