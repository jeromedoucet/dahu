<template>
  <div id="app">
    <nav-bar/>
    <div id='dahu-content'>
      <router-view></router-view> 
    </div>
  </div>
</template>

<script>

import NavBar from '@/components/NavBar.vue';

export default {
  name: 'app',
  components: {
    NavBar
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #cccccc;
  --darkest-color: #2B3E50;
  --dark-grey: #4E5D6C
}
</style>
