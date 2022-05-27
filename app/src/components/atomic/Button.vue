<template>
  <router-link v-if="to" :class="classes" :to="to">
    <slot />
  </router-link>
  <a v-else-if="href" :href="href" target="_blank" rel="noopener noreferrer" :class="classes">
    <slot />
  </a>
  <button v-else type="button" :class="classes">
    <slot />
  </button>
</template>

<script setup lang="ts">
import { RouteLocationRaw } from 'vue-router';

defineProps<{
  to?: RouteLocationRaw;
  href?: string;
}>();

const classes =
  'border-1 hover:border-gray-300 dark:border-gray-400 dark:hover:border-gray-500 rounded p-2 flex items-center justify-center cursor-pointer';
</script>
