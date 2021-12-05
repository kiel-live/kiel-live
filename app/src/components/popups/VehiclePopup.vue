<template>
  <div v-if="vehicle" class="flex flex-col">
    <span class="text-lg font-bold">{{ vehicle.name }}</span>
    <span class="pb-2 mb-2 border-b-1 dark:border-gray-500 text-sm">{{ vehicle.type }}</span>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from 'vue';
import { Marker } from '~/types';
import { vehicles } from '~/api';

export default defineComponent({
  name: 'VehiclePopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props) {
    const vehicle = computed(() => vehicles.value[props.marker.id]);

    return {
      vehicles,
      vehicle,
    };
  },
});
</script>
