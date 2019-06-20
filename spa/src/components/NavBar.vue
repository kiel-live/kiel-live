<template>
  <div class="navbar">
    <router-link :to="{ name: 'home' }" class="title">{{ title }}</router-link>
    <span v-if="!isConnected" class="offline"><i class="fas fa-signal"/>Offline</span>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'NavBar',
  data: () => ({
    title: process.env.VUE_APP_TITLE || 'OPNV Live',
  }),
  computed: {
    ...mapState([
      'isConnected',
    ]),
  },
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

  .offline {
    margin-left: auto;

    i {
      position: relative;
      margin-right: .5rem;

      &::after {
        position: absolute;
        top: 60%;
        left: 60%;
        width: 120%;
        height: .2rem;
        background: red;
        transform: translate(-50%, -50%) rotate(45deg);
        content: '';
      }
    }
  }
</style>
