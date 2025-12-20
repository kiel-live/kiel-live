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
import type { RouteLocationRaw } from 'vue-router';
import { computed } from 'vue';

const props = defineProps<{
  to?: RouteLocationRaw;
  href?: string;
  rounded?: boolean;
}>();

const classes = computed(() => [
  'border-2',
  'border-gray-200',
  'hover:border-gray-300',
  'dark:border-gray-400',
  'dark:hover:border-gray-500',
  'p-2',
  'flex',
  'items-center',
  'justify-center',
  'cursor-pointer',
  props.rounded ? 'rounded-full' : 'rounded',
]);
</script>
