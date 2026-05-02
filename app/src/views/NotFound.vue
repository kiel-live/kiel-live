<template>
  <div class="flex h-full w-full flex-col items-center justify-center p-4">
    <div class="m-auto flex flex-col">
      <p class="mb-8 text-2xl">{{ t('not_found_error') }}</p>
      <i18n-t keypath="not_found_back_to" tag="span" class="mx-auto">
        <template #link>
          <router-link class="text-blue-700 dark:text-blue-400" replace :to="{ name: 'home' }">
            {{ t('home') }}
          </router-link>
        </template>
      </i18n-t>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>
