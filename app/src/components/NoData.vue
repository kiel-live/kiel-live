<template>
  <div class="flex flex-col items-center flex-grow">
    <i-fa-solid-ban class="text-3xl mt-auto" />
    <div class="flex flex-col items-center my-2 text-lg"><slot /></div>
    <CustomButton class="mt-auto mb-4" :to="{ name: 'home' }" replace>
      <i-fa-close class="mr-2" /><span>Schließen</span>
    </CustomButton>
  </div>
</template>

<script setup lang="ts">
import CustomButton from '~/components/atomic/CustomButton.vue';
</script>
