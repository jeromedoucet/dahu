import Home from '@/components/Home.vue';
import JobCreation from '@/components/jobs/JobCreation.vue';
import JobDetail from '@/components/jobs/JobDetail.vue';
import Login from '@/components/Login.vue';
import RegistriesList from '@/components/registries/RegistriesList.vue'

import { isAuthenticated } from '@/services/user';

export const routes = [
  { path: '/login', component: Login },
  { path: '/', component: Home },
  { path: '/jobs/creation', component: JobCreation },
  { path: '/jobs/:jobId/pipeline', component: JobDetail },
  { path: '/registries', component: RegistriesList },
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
