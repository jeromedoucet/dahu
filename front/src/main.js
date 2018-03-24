import App from './App.vue'
import BootstrapVue from 'bootstrap-vue'
import 'bootswatch/dist/superhero/bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import './main.scss'
import Vue from 'vue/dist/vue.js';
import VueRouter from 'vue-router';
import { 
  routes,
  checkAuthenticationBeforeNavigation
} from './routes';

Vue.use(BootstrapVue);
Vue.config.productionTip = false;

Vue.use(VueRouter)
const router = new VueRouter({
  routes
})
router.beforeEach(checkAuthenticationBeforeNavigation);

new Vue({
  el: '#app',
  template: '<App/>',
  components: { App },
  router
});
