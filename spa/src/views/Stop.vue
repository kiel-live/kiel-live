<template>
  <div v-if="stop" class="stop">
    <router-link :to="{ name: 'home' }" class="back button"><i class="fas fa-angle-double-left"/></router-link>
    <h1 class="title">{{ stop.stopName }}</h1>
    <div class="arrivals">
      <div class="bus" v-for="bus in stop.actual" :key="bus.passageid">
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
      <div v-if="stop.actual.length == 0" class="no-data">
        <i class="fas fa-ban" />
        <p>Hier will gerade wohl kein Manni halten.</p>
      </div>
    </div>
  </div>
  <div v-else class="loading">
    <i class="fas fa-circle-notch fa-spin"></i>
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
        return 'sofort';
      }

      return `${minutes} Min`;
    },
  },
  mounted() {
    this.join();

    Api.on('connect', this.join);

    // wait for stop updates
    Api.on('stop', this.updateStop);
  },
  beforeDestroy() {
    Api.removeListener('connect', this.join);
    Api.removeListener('stop', this.updateStop);

    if (this.stop) {
      Api.emit('stop:leave', this.stop.stopShortName);
    }
  },
};
</script>

<style lang="scss" scoped>
  .stop {
    position: relative;
    display: flex;
    flex-flow: column;
    width: 100%;
    max-width: 40rem;
    margin: 0 auto;
    align-items: center;
    padding-top: 2rem;

    h1 {
      margin: 0;
      margin-bottom: 2rem;
    }

    .back {
      position: absolute;
      top: 2rem;
      left: .5rem;
    }
  }

  @media (min-width: 768px) {
    .stop {
      .back {
        left: 0;
      }
    }
  }

  .title {
    margin: 1rem 0;
    font-size: 1.8rem;
  }

  .arrivals {
    display: flex;
    flex-flow: column;
    width: 100%;
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

  .no-data {
    margin: auto;

    i {
      font-size: 4rem;
      margin-bottom: 1rem;
    }
  }

  .loading {
    margin: auto;
    font-size: 4rem;
  }
</style>
