<template>
  <div class="flex flex-col items-center flex-grow">
    <i-fa-solid-ban class="text-3xl mt-auto" />
    <p class="my-2 text-lg"><slot /></p>
    <Button to="home" class="mt-auto mb-4"><i-fa-close class="mr-2" /><span>Schließen</span> </Button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

import Button from '~/components/atomic/Button.vue';

export default defineComponent({
  name: 'NoData',
  components: { Button },
});
</script>
