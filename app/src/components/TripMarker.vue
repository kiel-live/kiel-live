<template>
  <div
    class="
      relative
      flex
      justify-center
      items-center
      mx-4
      h-12
      w-8
      after:(absolute
      left-1/2
      top-0
      h-full
      bg-gray-800
      w-0.5
      content-s)
    "
  >
    <div
      v-if="marker !== 'empty'"
      class="rounded-full h-4 w-4 flex items-center justify-center"
      :class="{
        'bg-gray-800': marker === 'dot',
        'bg-red-800': marker === 'pulse',
      }"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue';

export default defineComponent({
  name: 'TripMarker',

  props: {
    marker: {
      type: String as PropType<'empty' | 'pulse' | 'dot'>,
      required: true,
    },
  },
});
</script>
