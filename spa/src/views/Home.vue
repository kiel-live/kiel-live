<template>
  <div class="home">
    <div class="search">
      <TextInput class="searchText" :value="stopQuery" @input="searchStop" :debounce="true" :wait="400" placeholder="Haltestelle" autofocus />
      <div v-if="gpsSupport" class="button gps" :class="{ loading: gpsLoading }" @click="gps">
        <i v-if="gpsLoading" class="fas fa-circle-notch fa-spin" />
        <i v-else class="fas fa-crosshairs" />
      </div>
      <div class="map button" v-if="isTester" @click="$router.push({ name: 'map' })"><i class="fas fa-map-marked"/></div>
    </div>
    <div class="stops">
      <router-link v-for="stop in allStops" :key="stop.id" :to="{ name: 'stop', params: { stop: stop.id } }" class="stop" :class="{ favorite: stop.favorite }">
        <span class="name">{{ stop.name }}</span>
        <div class="details">
          <span v-if="stop.distance" class="direction">{{ stop.distance }}m</span>
          <i v-if="stop.favorite" class="icon fas fa-star"></i>
          <i  v-else-if="stop.distance" class="icon fas fa-crosshairs" />
          <i v-else class="icon fas fa-arrow-right"></i>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import Api from '@/api';
import TextInput from '@/components/TextInput.vue';

export default {
  name: 'home',
  components: {
    TextInput,
  },
  data() {
    return {
      gpsSupport: true,
      gpsLoading: false,
    };
  },
  computed: {
    ...mapGetters({
      allStops: 'stops/allStops',
    }),
    ...mapState({
      stopQuery: 'stops/stopQuery',
      isTester: 'isTester',
    }),
  },
  methods: {
    searchStop(query) {
      this.$store.commit('stops/setStopQuery', query);
      Api.emit('stop:search', query);
    },
    gps() {
      if (!this.gpsSupport || this.gpsLoading) { return; }

      this.gpsLoading = true;
      navigator.geolocation.getCurrentPosition((position) => {
        Api.emit('stop:nearby', {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        });
      }, (error) => {
        // eslint-disable-next-line no-console
        console.error(error);
        this.gpsSupport = false;
      });
    },
  },
  mounted() {
    this.gpsSupport = !!navigator.geolocation;

    Api.on('stop:search', (stops) => {
      this.$store.commit('stops/setStops', stops);
    });

    Api.on('stop:nearby', (stops) => {
      this.gpsLoading = false;
      this.$store.commit('stops/setGPSStops', stops);
    });
  },
};
</script>

<style lang="scss" scoped>
  .home {
    display: flex;
    flex-flow: column;
    margin: 0 auto;
    width: 100%;
    max-width: 40rem;
  }

  .search {
    display: flex;
    flex-flow: row;
    align-items: center;
    margin: 1rem .5rem;
    width: calc(100% - 1rem);

    .searchText {
      margin-right: 1rem;
    }

    .map {
      margin-left: 1rem;
    }
  }

  @media only screen and (min-width: 768px) {
    .search {
      width: 100%;
      margin: 1rem 0;
    }
  }

  .stops {
    display: flex;
    flex-flow: column;
  }

  .stop {
    display: flex;
    padding: 1rem;
    flex-flow: row;
    width: 100%;
    box-shadow: inset 0 -1px 0 0 rgba(100,121,143,0.122);
    text-align: left;
    text-decoration: none;
    color: #000;
    cursor: pointer;

    &:hover {
      -webkit-box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      box-shadow: inset 1px 0 0 #dadce0, inset -1px 0 0 #dadce0, 0 1px 2px 0 rgba(60,64,67,.3), 0 1px 3px 1px rgba(60,64,67,.15);
      z-index: 1;
    }

    &.favorite {
      i {
        color: gold;
      }
    }

    .details {
      display: flex;
      flex-flow: row;
      margin-left: auto;

      .direction {
        margin-right: .5rem;
      }

      .icon {
        width: 1rem;
        text-align: center;
      }
    }
  }
</style>
