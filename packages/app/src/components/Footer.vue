<template>
  <div v-if="visible" class="footer">
    <router-link :to="{name: 'about'}">
      {{ title }} - Made with &#10084; and Vue.js
    </router-link>
  </div>
</template>

<script>
import config from '@/libs/config';

export default {
  name: 'Footer',
  data: () => ({
    title: config('title', 'OPNV Live'),
  }),
  computed: {
    visible() {
      return this.$route.meta.disableFooter !== true;
    },
  },
};
</script>

<style lang="scss" scoped>
  .footer {
    display: flex;
    width: 100%;
    margin-top: auto;
    padding: .5rem 0;
    text-align: center;
    justify-content: center;

    a {
      color: #333;
      text-decoration: none;
    }
  }
</style>
