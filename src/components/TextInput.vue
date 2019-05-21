<template>
  <input @input="input" type="text" v-bind="$attrs" />
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
