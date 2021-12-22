<template>
  <div class="flex flex-col items-center">
    <i-fa-solid-ban class="text-3xl my-2" />
    <p class="my-2"><slot /></p>
    <router-link class="border-1 border-solid rounded p-2 flex justify-center" :to="{ name: 'home' }" replace>
      <i-fa-close class="mr-2" /><span>Schließen</span></router-link
    >
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'NoData',
});
</script>
