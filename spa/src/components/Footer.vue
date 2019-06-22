<template>
  <div class="footer">
    <router-link :to="{name: 'about'}">{{ title }} - Made with &#10084; and Vue.js</router-link>
  </div>
</template>

<script>
export default {
  name: 'Footer',
  data: () => ({
    title: process.env.VUE_APP_TITLE || 'OPNV Live',
  }),
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
      color: #000;
      text-decoration: none;
    }
  }
</style>
