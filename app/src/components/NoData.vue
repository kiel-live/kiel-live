<template>
  <div class="flex flex-col items-center flex-grow">
    <i-fa-solid-ban class="text-3xl mt-auto" />
    <div class="flex flex-col items-center my-2 text-lg"><slot /></div>
    <Button class="mt-auto mb-4" :to="{ name: 'home' }" replace>
      <i-fa-close class="mr-2" /><span>{{ t('close') }}</span>
    </Button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';

const { t } = useI18n();
</script>
