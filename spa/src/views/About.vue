<template>
  <div class="about">
    <h1 class="title">{{ title }}</h1>
    <p>Version: {{ version }}</p>
    <br />
    <div class="button" @click="setTester">{{ isTester ? 'Testmodus verlassen.' : 'Tester werden!' }}</div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import config from '@/libs/config';

export default {
  name: 'About',
  data: () => ({
    title: config('title', 'OPNV Live'),
    version: process.env.VUE_APP_VERSION,
  }),
  computed: {
    ...mapState([
      'isTester',
    ]),
  },
  methods: {
    setTester() {
      this.$store.commit('setTester', !this.isTester);
    },
  },
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
