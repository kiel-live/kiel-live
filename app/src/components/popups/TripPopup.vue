<template>
  <div class="flex flex-col">
    <span class="mb-2 border-b-1 text-lg">Trip</span>
    <span>{{ marker.id }}</span>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';
import { Marker } from '~/types';

export default defineComponent({
  name: 'TripPopup',

  props: {
    marker: {
      type: Object as PropType<Marker>,
      required: true,
    },
  },

  setup(props, { emit }) {},
});
</script>
