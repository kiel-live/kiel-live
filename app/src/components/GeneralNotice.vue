<template>
  <PopupNotice :show="!!alert" @close="close">
    <div v-if="alert" class="flex flex-col text-center gap-4">
      <span class="text-xl mb-2">{{ alert.title }}</span>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <span class="prose" v-html="body" />
    </div>

    <div class="flex flex-row w-full justify-center mt-8">
      <Button @click="close">{{ t('ok') }}</Button>
    </div>
  </PopupNotice>
</template>

<script setup lang="ts">
import { useStorage } from '@vueuse/core';
import { micromark } from 'micromark';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import PopupNotice from '~/components/PopupNotice.vue';
import { localStoragePrefix } from '~/compositions/useUserSettings';

const { t } = useI18n();

type Alert = {
  id: string;
  title: string;
  start?: Date;
  end?: Date;
  msg: string;
};

const alerts: Alert[] = [
  {
    id: 'kvg-rbl-umzug',
    title: '⚠️ Einschränkungen bei den Busdaten der KVG',
    start: new Date('2024-02-20 00:00:00'),
    end: new Date('2024-02-22 23:59:59'),
    msg: `Die KVG stellt vom **20. bis zum 22. Februar 2024** ihr rechnergestütztes Betriebsleitsystem (RBL) um,
 sodass es teilweise zu Einschränkungen bei den Busdaten kommen kann. [weitere Informationen](https://www.kvg-kiel.de/aktuelles/betriebliches/unser-rechnergestuetztes-betriebsleitsystem-rbl-zieht-um).`,
  },
];

const readAlerts = useStorage<string[]>(`${localStoragePrefix}.alerts`, []);

const alert = computed(() =>
  alerts.find((a) => {
    const currentlyActive = a.start && a.end ? a.start <= new Date() && a.end >= new Date() : true;
    return !readAlerts.value.includes(a.id) && currentlyActive;
  }),
);

function close() {
  if (!alert.value) {
    return;
  }
  readAlerts.value.push(alert.value.id);
}

const body = computed(() => (alert.value ? micromark(alert.value.msg.trim()) : null));
</script>
