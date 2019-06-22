<template>
  <div v-if="trip" class="trip">
    <h2 class="title">Abfahrten der {{ trip.routeName }} nach {{ trip.directionText }}</h2>
    <div class="stops">
      <template v-for="i in ['old', 'actual']">
        <div v-for="stop in trip[i]" :key="stop.tripId" class="stop" :class="i" @click="openStop(stop)">
          <div class="time">{{ stop.actualTime }}</div>
          <div class="marker">
            <i v-if="i == 'old'" class="fas fa-blank" />
            <i v-else class="fas fa-circle" />
          </div>
          <div class="name">{{ stop.stop.name }}</div>
        </div>
      </template>
    </div>
  </div>
  <div v-else class="trip loading">
    <i class="fas fa-circle-notch fa-spin"></i>
  </div>
</template>

<script>
import Api from '@/api';

export default {
  name: 'Trip',
  props: {
    id: {
      type: String,
      required: true,
    },
    vehicleId: {
      type: String,
      required: true,
    },
  },
  data: () => ({
    trip: null,
  }),
  methods: {
    updateTrip(trip) {
      this.trip = trip;
    },
    openStop(stop) {
      this.$router.push({ name: 'stop', params: { stop: stop.stop.shortName } });
    },
  },
  mounted() {
    Api.emit('trip', {
      tripId: this.id,
      vehicleId: this.vehicleId,  
    });

    // wait for stop updates
    Api.on('trip', this.updateTrip);
  },
};
</script>

<style lang="scss" scoped>
  .trip {
    display: flex;
    flex-flow: column;
    width: 100%;
  }

  .title {
    margin: 1rem 0;
    font-size: 1.8rem;
  }

  .stops {
    display: flex;
    flex-flow: column;
    width: 100%;
    max-width: 40rem;
    margin: 0 auto;
    align-items: center;
    padding-top: 1rem;
  }

  .stop {
    display: flex;
    padding: 0 1rem;
    flex-flow: row;
    width: 100%;
    // box-shadow: inset 0 -1px 0 0 rgba(100,121,143,0.122);
    text-align: left;
    align-items: center;
    cursor: pointer;

    &:hover {
      -webkit-box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      z-index: 1;
    }

    .time {
      width: 3rem;
    }

    .marker {
      position: relative;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0 1rem;
      height: 3rem;
      width: 2rem;

      &::after {
        position: absolute;
        left: calc(50% - 1px);
        top: 0;
        height: 100%;
        width: 2px;
        background: black;
        content: '';
      }
    }
  }
</style>
