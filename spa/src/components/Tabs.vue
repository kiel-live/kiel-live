<template>
  <div class="tabs">
    <div class="head">
      <div v-for="tab in tabs" :class="{ 'tab' : true, 'is-active': tab.isActive }" :key="tab.id" @click="selectTab(tab)">
        {{ tab.name }}
      </div>
    </div>

    <div class="body">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Tabs',

  data() {
    return {
      tabs: [],
      activeTab: parseInt((this.$route.hash || '#0').replace('#', ''), 10),
    };
  },

  watch: {
    activeTab(tab) {
      this.$router.replace({ hash: `#${tab}` });
    },
  },

  methods: {
    selectTab(selectedTab) {
      this.tabs.forEach((tab) => {
        tab.isActive = false;
      });
      selectedTab.isActive = true;
      this.activeTab = selectedTab.name.toLowerCase();
    },
  },

  created() {
    this.tabs = this.$children;
  },
};
</script>

<style lang="scss" scoped>
.tabs {
  display: flex;
  flex-flow: column;
  width: 100%;

  .head {
    display: flex;
    width: 100%;
    border-bottom: 1px solid #eee;

    .tab {
      flex-grow: 1;
      padding: 1rem;
      cursor: pointer;

      &.is-active {
        border-bottom: 1px solid #aaa;
      }
    }
  }
}
</style>
