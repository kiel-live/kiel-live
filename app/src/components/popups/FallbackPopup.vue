<template>
  <div v-if="stop" class="flex flex-col min-h-0 flex-grow">
    <div class="flex flex-row pb-2 mb-2 border-b-1 dark:border-dark-100 items-center">
      <i-material-symbols-label-sharp />
      <h1 class="text-lg ml-2">{{ stop.name }}</h1>
      <Button
        v-if="isFavorite(stop)"
        class="text-yellow-300 ml-auto border-0"
        :title="t('remove_favorite')"
        @click="removeFavorite(stop)"
      >
        <i-ph-star-fill />
      </Button>
      <Button v-else class="ml-auto border-0" :title="t('add_favorite')" @click="addFavorite(stop)">
        <i-ph-star-bold />
      </Button>
    </div>
  </div>
  <NoData v-else>
    {{ t('this_stop_probably_does_not_exist') }}
    <Button
      v-if="isFavorite(marker)"
      class="mt-2"
      @click="
        () => {
          removeFavorite(marker);
          $router.replace({ name: 'home' });
        }
      "
    >
      <i-ph-star-fill class="mr-2 text-yellow-300" /><span>{{ t('remove_favorite') }}</span>
    </Button>
  </NoData>
</template>

<script setup lang="ts">
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import { stops } from '~/api';
import { Marker } from '~/api/types';
import Button from '~/components/atomic/Button.vue';
import NoData from '~/components/NoData.vue';
import { useFavorites } from '~/compositions/useFavorites';

const props = defineProps<{
  marker: Marker;
}>();

const { addFavorite, removeFavorite, isFavorite } = useFavorites();
const { t } = useI18n();

const marker = toRef(props, 'marker');
const stop = computed(() => stops.value[props.marker.id]);
</script>
