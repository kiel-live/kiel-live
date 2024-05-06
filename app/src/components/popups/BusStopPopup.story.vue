<template>
  <Story>
    <DetailsPopup is-open>
      <BusStopPopup :marker="marker" />
    </DetailsPopup>
  </Story>
</template>

<script lang="ts" setup>
import { Marker } from '~/api/types';
import DetailsPopup from '~/components/DetailsPopup.vue';

import BusStopPopup from './BusStopPopup.vue';

const marker: Marker = {
  id: '123',
  type: 'bus-stop',
};
</script>
