<template>
  <div v-if="driving" class="trip">
    <div class="header">
      <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left"/></a>
      <h1 class="title">{{ trip.routeName }} nach {{ trip.directionText }}</h1>
      <div class="map button" @click="$router.push({ name: 'mapTrip', params: { vehicle: vehicleId }})"><i class="fas fa-map-marked"/></div>
    </div>
    <div class="stops">
      <template v-for="i in ['old', 'actual']">
        <div v-for="(stop, index) in trip[i]" :key="stop.tripId" class="stop" :class="i" @click="openStop(stop)">
          <div class="time">{{ stop.actualTime }}</div>
          <div class="marker">

            <div v-if="i === 'actual' && index === 0" class="vehicle" :class="{ driving: (stop.status === 'PREDICTED') }">
              <div class="ringring"></div>
            </div>

            <i v-if="i === 'old'" class="fas fa-blank" />
            <i v-else class="fas fa-circle" />
          </div>
          <div class="name">{{ stop.stop.name }}</div>
        </div>
      </template>
    </div>
  </div>
  <div v-else-if="!tripId || !vehicleId" class="trip">
    <div class="notfound">
      <i class="fas fa-ban" />
      <p>Zu dieser Tour gibt es anscheinend keine live Daten.</p>
      <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left" />Zurück</a>
    </div>
  </div>
  <div v-else-if="!trip" class="trip loading">
    <i class="fas fa-circle-notch fa-spin"></i>
  </div>
  <div v-else class="trip">
    <div class="ended">
      <i class="fas fa-ban" />
      <p>Diese Tour ist wohl schon zu Ende.</p>
      <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left" />Zurück</a>
    </div>
  </div>
</template>

<script>
import Api from '@/api';

export default {
  name: 'Trip',
  data: () => ({
    trip: null,
    ids: {},
  }),
  computed: {
    tripId() {
      return this.$route.params.trip;
    },
    vehicleId() {
      return this.$route.params.vehicle;
    },
    driving() {
      return this.trip && (this.trip.old.length > 0 || this.trip.actual.length > 0);
    },
  },
  methods: {
    load() {
      if (!this.tripId || !this.vehicleId) {
        return;
      }

      this.join();
      Api.on('connect', this.join);
      // wait for trip updates
      Api.on('trip', this.updateTrip);
    },
    unload() {
      Api.removeListener('connect', this.join);
      Api.removeListener('trip', this.updateTrip);

      if (this.trip) {
        Api.emit('trip:leave', this.ids);
      }

      this.stop = null;
    },
    reload() {
      this.unload();
      this.load();
    },
    join() {
      this.ids = {
        tripId: this.tripId,
        vehicleId: this.vehicleId,
      };

      Api.emit('trip:join', this.ids);
    },
    updateTrip(trip) {
      this.trip = trip;
    },
    openStop(stop) {
      this.$router.push({ name: 'stop', params: { stop: stop.stop.shortName } });
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
  .trip {
    position: relative;
    display: flex;
    flex-flow: column;
    flex-grow: 1;
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
      justify-content: space-between;
      align-items: center;
    }

    .back {
      margin-right: 1rem;
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
    margin-bottom: 1rem;
    font-size: 1.8rem;
  }

  .stops {
    display: flex;
    flex-flow: column;
    width: 100%;
    max-width: 40rem;
    margin: 0 auto;
    align-items: center;
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
        left: calc(50% - .05rem);
        top: 0;
        height: 100%;
        width: .1rem;
        background: #2e2e2e;
        content: '';
      }

      i {
        z-index: 1;
      }
    }

    .vehicle {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 2;

      &.driving {
        top: -.25rem;
      }

      &::before {
        display: block;
        width: 1rem;
        height: 1rem;
        border-radius: 50%;
        background-color: #c80000;
        content: "";
      }

      .ringring {
        border: 3px solid #c80000;
        border-radius: 30px;
        height: 2rem;
        width: 2rem;
        position: absolute;
        left: calc(50% - 1rem);
        top: calc(50% - 1rem);
        transform: translate(-50%, -50%);
        animation: pulsate 1.5s ease-out;
        animation-iteration-count: infinite;
        opacity: 0.0
      }
    }
  }

  .loading {
    margin: auto;
    font-size: 4rem;
  }

  .trip .ended,
  .trip .notfound {
    margin: auto;

    > i {
      font-size: 4rem;
      margin-bottom: 1rem;
    }

    .back {
      margin-top: 1rem;

      i {
        margin-right: .5rem;
      }
    }
  }

  @keyframes pulsate {
    0% {-webkit-transform: scale(0.1, 0.1); opacity: 0.0;}
    50% {opacity: 1.0;}
    100% {-webkit-transform: scale(1.2, 1.2); opacity: 0.0;}
  }
</style>
