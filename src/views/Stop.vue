<template>
  <div v-if="stop" class="stop">
    <router-link :to="{ name: 'home' }">Back &lt;&lt;</router-link>
    <h2 class="title">{{ stop.stopName }} ({{ stopId }})</h2>
    <div class="arrivals">
      <div class="bus" v-for="bus in stop.actual" :key="bus.tripId">
        <div class="icon">
          <i v-if="route(bus.routeId).routeType === 'bus'" class="fas fa-bus"></i>
          <i v-if="route(bus.routeId).routeType === 'ferry'" class="fas fa-bus"></i>
        </div>
        <div class="line">{{ route(bus.routeId).shortName }}</div>
        <div class="direction">{{ bus.direction }}</div>
        <div class="eta">{{ eta(bus) }}</div>
        <div class="status">
          <i v-if="bus.status === 'STOPPING'" class="fas fa-hand-paper"></i>
          <i v-if="bus.status === 'PLANNED'" class="fas fa-clock"></i>
          <i v-if="bus.status === 'PREDICTED'" class="fas fa-running"></i>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Api from '@/api';

export default {
  name: 'stop',
  data() {
    return {
      stop: null,
    };
  },
  computed: {
    stopId() {
      return this.$route.params.stop;
    },
  },
  methods: {
    join() {
      // request server to join stop room
      console.log('joined', this.stopId);
      Api.emit('stop:join', this.stopId);
    },
    updateStop(stop) {
      this.stop = stop;
    },
    route(routeId) {
      for (let i = 0; i < this.stop.routes.length; i += 1) {
        if (this.stop.routes[i].id === routeId) {
          return this.stop.routes[i];
        }
      }

      return null;
    },
    eta(bus) {
      const minutes = Math.round(bus.actualRelativeTime / 60);

      if (bus.status === 'STOPPING') {
        return 'hält';
      }

      if (bus.status === 'PLANNED') {
        return bus.plannedTime;
      }

      if (minutes < 1) {
        return `sofort (${bus.actualRelativeTime} Sek)`;
      }

      return `${minutes} Min`;
    },
  },
  mounted() {
    console.log('toll');
    this.join();

    Api.on('connect', this.join);

    // wait for stop updates
    Api.on('stop', this.updateStop);
  },
  beforeDestroy() {
    Api.removeListener('connect', this.join);
    Api.removeListener('stop', this.updateStop);

    if (this.stop) {
      console.log('left', this.stop.stopShortName);
      Api.emit('stop:leave', this.stop.stopShortName);
    }
  },
};
</script>

<style lang="scss">
  .stop {
    display: flex;
    flex-flow: column;
    width: 100%;
    align-items: center;
  }

  .title {
    font-size: 1.8rem;
    margin-bottom: 1rem;
  }

  .arrivals {
    display: flex;
    flex-flow: column;
    width: 100%;
    max-width: 40rem;
  }

  .bus {
    display: flex;
    padding: 1rem;
    flex-flow: row;
    width: 100%;
    box-shadow: inset 0 -1px 0 0 rgba(100,121,143,0.122);
    text-align: left;
    cursor: pointer;

    &:hover {
      -webkit-box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      z-index: 1;
    }

    * {
      display: flex;
      align-items: center;
    }

    .icon {
      margin-right: .5rem;
    }

    .line {
      width: 1.5rem;
    }

    .direction {
      margin-left: 1rem;
      flex-grow: 1;
    }

    .eta {
      width: 20%;
      justify-content: flex-end;
    }

    .status {
      width: 10%;
      justify-content: center;
      font-size: 1.5rem;
    }
  }
</style>
