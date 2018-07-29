<template>
  <b-navbar id='dahu-navbar' toggleable="md" type="dark" variant="dark">
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
    <b-navbar-brand href="#/">Dahu</b-navbar-brand> 
    <b-navbar-nav v-if="isAuthenticated()">
      <b-nav-item id="registries-nav-link" href="#/registries">Images registries</b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav class="ml-auto" v-if="isAuthenticated()">
      <b-nav-form right>
        <b-button size="sm" class="my-2 my-sm-0" type="button" v-on:click="disconnect">Disconnect</b-button>
      </b-nav-form>
    </b-navbar-nav>
  </b-navbar>
</template>

<script>
import { logout, isAuthenticated } from '@/services/user';
export default {
  name: 'NavBar',
  props: {
  },
  methods: {
    disconnect () {
      logout();
      this.$router.push('/login');
    },
    isAuthenticated() {
      return isAuthenticated();
    }
  }
}
</script>

<style scoped>
#dahu-navbar {
  margin-bottom: 3rem;
}
</style>
