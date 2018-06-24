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

input[type='text'] {
  background: transparent;
  border: none;
  border-bottom: 1px solid #000000;
  font-size: 16px;
  font-family: 'HK Grotesk';
}

.item-input {
  border-bottom: 0px !important;
  margin-bottom: 8px;
}

.green {
  color: #5cb85c;
}

.pretty {
  /* Overwrite pretty checkmark CSS margin-right */
  margin-right: 0;
}

.strike {
  text-decoration: line-through;
  color: #5cb85c;
  -webkit-transition: all 0.2s ease;
  -moz-transition: all 0.2s ease;
  -o-transition: all 0.2s ease;
  -ms-transition: all 0.2s ease;
  transition: all 0.2s ease;
}

.round-btn {
  border-radius: 100%;
}
</style>
