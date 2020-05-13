<template>
  <div v-if="visible" class="modal" aria-modal="true">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head"><p class="modal-card-title">Neue Funktionen</p></header>
      <section class="modal-card-body">
        <div v-for="change in changes" :key="change.timestamp" class="change-block">
          <p class="change-title">{{ change.timestamp }}</p>
          <span class="change-content" v-html="change.content" />
        </div>
      </section>
      <footer class="modal-card-foot">
        <div class="button" @click="setVersion">Okay<i class="fas fa-rocket icon"/></div>
      </footer>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import changelog from '@/changelog.json';

export default {
  name: 'Changelog',

  computed: {
    ...mapState([
      'version',
    ]),
    changes() {
      return changelog.changes;
    },
    visible() {
      return this.version !== changelog.version;
    },
  },

  methods: {
    setVersion() {
      this.$store.commit('setVersion', changelog.version);
    },
  },
};
</script>

<style lang="scss" scoped>
  .modal {
    display: flex;
    align-items: center;
    flex-direction: column;
    justify-content: center;
    overflow: hidden;
    position: fixed;
    z-index: 10000;
    bottom: 0;
    left: 0;
    right: 0;
    top: 0;

    &-background {
      background-color: hsla(0,0%,4%,.86);
      bottom: 0;
      left: 0;
      position: absolute;
      right: 0;
      top: 0;
    }

    &-card {
      display: flex;
      width: calc(100% - 2rem);
      flex-grow: 1;
      flex-flow: column;
      margin: 1rem;
      max-width: 30rem;
      max-height: 30rem;

      &-foot, &-head {
        align-items: center;
        background-color: #f5f5f5;
        display: flex;
        flex-shrink: 0;
        padding: 20px;
        position: relative;
      }

      &-head {
        border-bottom: 1px solid #dbdbdb;
        border-top-left-radius: 6px;
        border-top-right-radius: 6px;
        justify-content: flex-start;
      }

      &-title {
        color: #363636;
        flex-grow: 1;
        flex-shrink: 0;
        font-size: 1.5rem;
        line-height: 1;
      }

      &-body {
        background-color: #fff;
        flex-grow: 1;
        flex-shrink: 1;
        overflow: auto;
        padding: 20px;
        position: relative;
        text-align: left;

        .change-block:not(:last-child) {
          margin-bottom: 1.5rem;
        }

        .change-title {
          position: relative;
          font-size: 1.2rem;
          margin-bottom: 1rem;

          &::after {
            content: " ";
            position: absolute;
            background: #dbdbdb;
            left: 0;
            bottom: -0.5rem;
            width: 100%;
            height: 1px;
          }
        }

        .change-content {
          line-height: 1.5rem;
        }

        ::v-deep b {
          font-weight: bold;
        }
      }

      &-foot {
        border-bottom-left-radius: 6px;
        border-bottom-right-radius: 6px;
        border-top: 1px solid #dbdbdb;
        justify-content: center;

        .button {
          background-color: #363636;
          border-color: #363636;
          color: #fff;
        }

        .icon {
            margin-left: .5rem;
        }
      }
    }
  }
</style>
