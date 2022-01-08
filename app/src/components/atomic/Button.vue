<template>
  <div
    class="border-1 hover:border-gray-300 dark:border-gray-400 dark:hover:border-gray-500 rounded p-2 flex justify-center cursor-pointer"
  >
    <router-link v-if="to" class="flex items-center" :to="{ name: to }" replace>
      <slot>
        <span>{{ text }}</span>
      </slot>
    </router-link>
    <button v-else type="button">
      <slot>
        <span>{{ text }}</span>
      </slot>
    </button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'Button',

  props: {
    to: {
      type: String,
      default: null,
    },

    text: {
      type: String,
      default: null,
    },
  },
});
</script>
