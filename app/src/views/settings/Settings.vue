<template>
  <div class="flex flex-col min-h-0 flex-grow overflow-y-auto p-2">
    <h1 class="mb-2 text-xl font-bold">{{ t('settings') }}</h1>

    <div class="flex">
      <div class="flex flex-col">
        <span>Litemode aktivieren</span>
        <span class="text-sm dark:text-gray-400"
          >Wenn du den Litemode aktiviert hast wird die Karte nicht angezeigt, sodass du auch mit geringem Datenvolumen
          die App nutzen können solltest.</span
        >
      </div>
      <Checkbox />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Checkbox from '~/components/atomic/Checkbox.vue';

const { t } = useI18n();
</script>
