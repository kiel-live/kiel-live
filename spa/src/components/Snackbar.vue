<template>
  <div class="snackbar" :class="{ show }">
    <slot>
      {{ text }}
      <div v-if="action" class="action" @click="click">{{ action }}</div>
    </slot>
  </div>
</template>

<script>
export default {
  name: 'Snackbar',
  props: {
    text: {
      required: true,
    },
    show: {
      default: false,
    },
    action: {
      default: null,
    },
  },
  methods: {
    click() {
      this.$emit('click');
    },
  },
};
</script>

<style lang="scss" scoped>
.snackbar {
  position: fixed;
  display: flex;
  flex-flow: row;
  align-items: center;
  width: 100%;
  left: 50%;
  bottom: -4rem;
  background-color: #323232;
  color: #fff;
  text-align: center;
  border-radius: 2px;
  padding: 1rem;
  z-index: 1;
  transform: translateX(-50%);
  transition: all .5s;

  /* Show the snackbar when clicking on a button (class added with JavaScript) */
  &.show {
    bottom: 0;
  }

  .action {
    margin-left: 1rem;
    cursor: pointer;
    text-transform: uppercase;
    color: #4CAF50;
  }
}

@media (min-width: 768px) {
  .snackbar {
    width: auto;
    min-width: 250px;

    &.show {
      bottom: 2rem;
    }

    .action {
      margin-left: 2rem;
    }
  }
}
</style>
