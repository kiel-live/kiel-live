<template>
  <div class="flex flex-col h-full w-full items-center justify-center p-4">
    <div class="flex flex-col m-auto">
      <p class="text-2xl mb-8">{{ t('404_error') }}</p>
      <i18n-t keypath="404_back_to" tag="span" class="mx-auto">
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
