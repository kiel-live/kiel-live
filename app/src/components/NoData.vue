<template>
  <div class="flex grow flex-col items-center">
    <i-ph-empty class="mt-auto text-3xl" />
    <div class="my-2 flex flex-col items-center text-lg"><slot /></div>
    <Button class="mt-auto mb-4" :to="{ name: 'home' }" replace>
      <i-ph-x-bold class="mr-2 font-bold" /><span>{{ t('close') }}</span>
    </Button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';

const { t } = useI18n();
</script>
