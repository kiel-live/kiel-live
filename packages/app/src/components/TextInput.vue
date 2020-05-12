<template>
  <div :class="{ input: true, typing }">
    <input @input="input" type="text" v-model="innerValue" v-bind="$attrs" />
  </div>
</template>

<script>
export default {
  name: 'TextInput',
  inheritAttrs: false,
  data() {
    return {
      timeout: null,
      typing: false,
      innerValue: null,
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
    value: {
      type: String,
      default: null,
    },
  },
  methods: {
    input(event) {
      this.typing = true;

      if (!this.debounce) {
        this.$emit('input', event.target.value);
        return;
      }

      if (this.timeout) {
        clearTimeout(this.timeout);
      }

      this.timeout = setTimeout(() => {
        this.typing = false;
        this.$emit('input', event.target.value);
      }, this.wait);
    },
  },
  mounted() {
    this.innerValue = this.value;
  },
};
</script>

<style lang="scss" scoped>
  @keyframes spinAround {
    0% {
      transform: rotate(0deg)
    }
    100% {
      transform: rotate(360deg)
    }
  }

  .input {
    position: relative;
    width: 100%;
    margin: 8px 0;

    input {
      display: inline-flex;
      width: 100%;
      max-width: 100%;
      border: 1px solid #ccc;
      box-shadow: inset 0 1px 3px #ddd;
      font-size: 1rem;
      border-radius: .2rem;
      padding: .5rem .5rem;

      &:focus {
        box-shadow: inset 0 1px 1px rgba(0, 0, 0, .075), 0 0 8px rgba(102, 175, 233, .6);
      }
    }

    &.typing {
      input {
        padding-right: 2.25em;
      }

      &::after {
        position: absolute !important;
        right: .625em;
        top: .625em;
        z-index: 4;

        animation: spinAround .5s infinite linear;
        border: 2px solid #dbdbdb;
        border-radius: 290486px;
        border-right-color: transparent;
        border-top-color: transparent;
        content: "";
        display: block;
        height: 1em;
        position: relative;
        width: 1em;
      }
    }
  }
</style>
