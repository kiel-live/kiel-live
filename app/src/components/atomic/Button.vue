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
  'border',
  'border-gray-200',
  'hover:border-gray-300',
  'dark:border-neutral-950',
  'dark:hover:border-neutral-700',
  'hover:bg-gray-50',
  'dark:hover:bg-neutral-900',
  'px-4',
  'py-2.5',
  'flex',
  'items-center',
  'justify-center',
  'cursor-pointer',
  'transition-all',
  'font-medium',
  props.rounded ? 'rounded-full' : 'rounded-xl',
]);
</script>
