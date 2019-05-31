<template>
  <div class="navbar">
    <router-link :to="{ name: 'home' }" class="title">Live OPNV</router-link>
  </div>
</template>

<script>
export default {
  name: 'NavBar',
};
</script>

<style lang="scss" scoped>
  .navbar {
    display: flex;
    width: 100%;
    padding: 1rem;
    box-shadow: 0 1px 8px -1px rgba(0, 0, 0, .5);
  }

  .right {
    margin-left: auto;
  }

  .title {
    color: #333;
    text-decoration: none;
    font-size: 1.1rem;
  }
</style>
