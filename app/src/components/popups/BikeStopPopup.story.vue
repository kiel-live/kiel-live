<template>
  <Story>
    <DetailsPopup is-open>
      <BikeStopPopup />
    </DetailsPopup>
  </Story>
</template>

<script lang="ts" setup>
import DetailsPopup from '~/components/DetailsPopup.vue';

import BikeStopPopup from './BikeStopPopup.vue';
</script>
