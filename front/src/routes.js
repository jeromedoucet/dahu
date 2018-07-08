import Home from '@/components/Home.vue';
import JobCreation from '@/components/jobs/JobCreation.vue';
import Login from '@/components/Login.vue';

import { isAuthenticated } from '@/services/user';

export const routes = [
  { path: '/login', component: Login },
  { path: '/', component: Home },
  { path: '/jobs/creation', component: JobCreation }
]

export function checkAuthenticationBeforeNavigation(to, from, next) {
  if(isAuthenticated()) {
    if(to.path === '/login') {
      next('/');
    } else {
      next();
    }
  } else {
    if(to.path === '/login') {
      next();
    } else {
      next('/login');
    }
  }
}
