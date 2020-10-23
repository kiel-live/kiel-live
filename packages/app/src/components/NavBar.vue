<template>
  <div class="top">
    <div class="navbar">
      <div class="inner">
        <router-link :to="{ name: 'home' }" class="title">
          {{ title }}
        </router-link>
        <span v-if="!isConnected" class="offline"><i class="fas fa-signal" />Offline</span>
      </div>
    </div>
    <div class="placeholder" />
  </div>
</template>

<script>
import { mapState } from 'vuex';
import config from '@/libs/config';

export default {
  name: 'NavBar',
  data: () => ({
    title: config('title', 'OPNV Live'),
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
    position: fixed;
    display: flex;
    width: 100%;
    padding: 1rem;
    box-shadow: 0 1px 8px -1px rgba(0, 0, 0, .5);
    z-index: 1000;
    background: #fff;

    .inner {
      position: relative;
      display: flex;
      flex-flow: row;
      width: 100%;
      max-width: 40rem;
      margin: 0 auto;
      padding: 0 0.5rem;
    }
  }

  .top {
    position: relative;
  }

  .placeholder {
    display: block;
    height: calc(2rem + (1rem + 1px)); // padding + text size
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
