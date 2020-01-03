<template>
  <div class="about">
    <h1 class="title">{{ title }}</h1>
    <p>Version: {{ version }}</p>
  </div>
</template>

<script>
import config from '@/libs/config';

export default {
  name: 'About',
  data: () => ({
    title: config('title', 'OPNV Live'),
    version: process.env.VUE_APP_VERSION,
  }),
};
</script>

<style lang="scss" scoped>
  .about {
    margin: auto 0;
  }

  .title {
    margin-bottom: 1rem;
    font-size: 1.8rem;
  }
</style>
