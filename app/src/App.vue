<template>
  <div class="app flex flex-col m-auto w-full h-full bg-gray-100">
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'App',

  components: {},
});
</script>

<style>
html,
body,
#app {
  width: 100%;
  height: 100%;
}

body {
  /* disable android pull to refresh feature */
  overflow-y: hidden;
}
</style>

<style scoped>
.app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
