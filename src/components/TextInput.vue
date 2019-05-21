<template>
  <input class="input" @input="input" type="text" v-bind="$attrs" />
</template>

<script>
export default {
  name: 'TextInput',
  data() {
    return {
      stops: null,
      timeout: null,
    };
  },
  props: {
    debounce: {
      type: Boolean,
      default: false,
    },
    wait: {
      type: Number,
      default: 1000,
    },
  },
  methods: {
    input(event) {
      const { value } = event.target;

      if (!this.debounce) {
        this.$emit('input', value);
        return;
      }

      if (this.timeout) {
        clearTimeout(this.timeout);
      }

      this.timeout = setTimeout(() => {
        this.$emit('input', value);
      }, this.wait);
    },
  },
};
</script>

<style lang="scss" scoped>
  .input {
    display: inline-block;
    width: 100%;
    margin: 8px 0;
    border: 1px solid #ccc;
    box-shadow: inset 0 1px 3px #ddd;
    font-size: 1rem;
    border-radius: .2rem;
    padding: .5rem .5rem;
  }

  .input:focus {
    box-shadow: inset 0 1px 1px rgba(0, 0, 0, .075), 0 0 8px rgba(102, 175, 233, .6);
  }
</style>
