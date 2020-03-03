<template>
  <div class="map-container">
    <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left"/></a>
    <p>Wird bald zur Verfügung stehen ;-D</p>
  </div>
</template>

<script>
export default {
  name: 'osmap',
  props: {
    focusVehicle: {
      type: String,
      default: null,
    },
    focusStop: {
      type: String,
      default: null,
    },
  },
};
</script>

<style lang="scss" scoped>
  .map-container {
    position: relative;
    display: flex;
    width: 100%;
    height: 100%;
    justify-content: center;
    align-items: center;

    .back {
      margin-right: 1rem;
    }
  }
</style>
