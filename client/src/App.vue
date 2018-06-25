<template>
  <div id="app">
    <Navbar/>
    <router-view/>
  </div>
</template>

<script>
import Navbar from './components/Navbar'

export default {
  name: 'App',
  components: {
    Navbar
  }
}
</script>

<style>
@font-face {
  font-family: 'HK Grotesk';
  src: url('assets/hk-grotesk/HKGrotesk-Medium.otf');
}

#app {
  text-align: center;
  font-family: 'HK Grotesk';
}
</style>
