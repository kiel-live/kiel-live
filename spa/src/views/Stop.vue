<template>
  <div v-if="stop" class="stop">
    <div class="header">
      <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left"/></a>
      <h1 class="title">{{ stop.stopName }}</h1>
      <div v-if="isFavorite" class="favorite gold button" @click="removeFavoriteStop"><i class="fas fa-star"/></div>
      <div v-else class="favorite button" @click="addFavoriteStop"><i class="far fa-star"/></div>
      <div class="map button" @click="$router.push({ name: 'mapStop', params: { stop: stopId }})"><i class="fas fa-map-marked"/></div>
    </div>
    <div v-if="alerts && alerts.length > 0" class="alerts">
      <div v-for="alert in alerts" :key="alert.title" class="alert">
        <i class="fas fa-exclamation-triangle icon" />
        <div class="content">{{ alert.title }}</div>
      </div>
    </div>
    <div class="arrivals">
      <div class="bus" v-for="bus in arrivals" :key="bus.passageid" @click="openTrip(bus)">
        <div class="icon">
          <i v-if="route(bus.routeId).routeType === 'bus'" class="fas fa-bus"></i>
          <i v-if="route(bus.routeId).routeType === 'ferry'" class="fas fa-bus"></i>
        </div>
        <div class="line">{{ route(bus.routeId).name }}</div>
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
import { orderBy } from 'lodash';
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
    isFavorite() {
      return this.$store.state.favoriteStops.filter((i) => i.id === this.stopId).length === 1;
    },
    arrivals() {
      if (!this.stop) {
        return [];
      }
      return orderBy(this.stop.actual, (stop) => {
        if (stop.status === 'STOPPING') {
          return 0;
        }

        if (stop.actualRelativeTime) {
          // if eta is delayed set it to 1 to be greater than STOPPING arrivals
          return Math.max(stop.actualRelativeTime, 1);
        }

        return stop.plannedTime;
      });
    },
    alerts() {
      if (!this.stop) {
        return [];
      }

      let alerts = this.stop.generalAlerts;

      this.stop.routes.forEach((route) => {
        alerts = [...alerts, ...route.alerts];
      });

      return alerts;
    },
  },
  methods: {
    load() {
      this.join();
      Api.on('connect', this.join);
      // wait for stop updates
      Api.on('stop', this.updateStop);
    },
    unload() {
      Api.removeListener('connect', this.join);
      Api.removeListener('stop', this.updateStop);

      if (this.stop) {
        Api.emit('stop:leave', this.stop.stopShortName);
      }

      this.stop = null;
    },
    reload() {
      this.unload();
      this.load();
    },
    join() {
      // request server to join stop room
      Api.emit('stop:join', this.stopId);
    },
    updateStop(stop) {
      this.stop = stop;

      if (stop) {
        this.$store.commit('setTitle', stop.stopName);
      }
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
    addFavoriteStop() {
      if (this.stop && this.stop.stopName) {
        this.$store.commit('addFavoriteStop', { id: this.stopId, name: this.stop.stopName });
      }
    },
    removeFavoriteStop() {
      this.$store.commit('removeFavoriteStop', this.stopId);
    },
    openTrip(bus) {
      this.$router.push({ name: 'trip', params: { trip: bus.tripId, vehicle: bus.vehicleId } });
    },
  },
  mounted() {
    this.load();
  },
  beforeDestroy() {
    this.unload();
  },
  beforeRouteUpdate(to, from, next) {
    next();
    this.reload();
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
    padding-top: 1rem;

    h1 {
      margin: 0;
    }

    .header {
      position: relative;
      width: calc(100% - 1rem);
      display: flex;
      margin: 1rem .5rem;
      margin-bottom: 2rem;
      align-items: center;
      justify-content: space-between;
    }

    .back {
      margin-right: .5rem;
    }

    .title {
      flex-grow: 2;
    }

    .favorite {
      margin-left: 1rem;

      &.gold {
        color: gold;
      }
    }

    .map {
      margin-left: 1rem;
    }
  }

  @media (min-width: 768px) {
    .stop {
      .header {
        width: 100%;
        margin-left: 0;
        margin-right: 0;
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
      margin-left: 1.5rem;
      flex-grow: 1;
    }

    .eta {
      width: 20%;
      justify-content: flex-end;
    }

    .status {
      width: 10%;
      justify-content: flex-end;
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

  .alerts {
    display: flex;
    width: 100%;
    padding-bottom: 2rem;
    flex-flow: column;

    .alert {
      display: flex;
      flex-flow: row;
      color: #fff;
      border-radius: 4px;
      padding: 1.25rem 2.5rem 1.25rem 1.5rem;
      position: relative;
      width: 100%;
      margin-bottom: .5rem;

      background-color: #ff3860;
      border-color: rgba(0,0,0,0.12) !important;

      .icon {
        margin-right: 1rem;
      }
    }
  }
</style>
