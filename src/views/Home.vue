<template>
  <div class="home">
    <TextInput class="search" @input="searchStop" :debounce="true" :wait="400" placeholder="Haltestelle" autofocus />
    <div v-if="stops" class="stops">
      <router-link v-for="stop in stops" :key="stop.id" :to="{ name: 'stop', params: { stop: stop.id } }" class="stop">
        <span class="name">{{ unEscapeHtml(stop.name) }}</span>
        <i class="icon fas fa-arrow-right"></i>
      </router-link>
    </div>
  </div>
</template>

<script>
import Api from '@/api';

import TextInput from '@/components/TextInput.vue';

export default {
  name: 'home',
  components: {
    TextInput,
  },
  data() {
    return {
      stops: null,
    };
  },
  methods: {
    searchStop(query) {
      Api.emit('stop:search', query);
    },
    unEscapeHtml(unsafe) {
      return unsafe
        .replace('&amp;', /&/g)
        .replace('&lt;', /</g)
        .replace('&gt;', />/g)
        .replace('&quot;', /"/g)
        .replace('&#039;', /'/g)
        .replace('&auml;', 'ä')
        .replace('&Auml;', 'Ä')
        .replace('&ouml;', 'ö')
        .replace('&Ouml;', 'Ö')
        .replace('&uuml;', 'ü')
        .replace('&Uuml;', 'Ü')
        .replace('&szlig;', 'ß');
    },
  },
  mounted() {
    Api.on('stop:search', (stops) => {
      this.stops = [];

      const regexp = RegExp('<li stop="(.*?)">(.*?)</li>', 'g');
      let matches = regexp.exec(stops);
      while (matches) {
        this.stops.push({
          id: matches[1],
          name: matches[2],
        });
        matches = regexp.exec(stops);
      }
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
    margin: 0 .5rem;
    margin-top: 1rem;
    width: calc(100% - 1rem);
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

    .icon {
      margin-left: auto;
    }
  }
</style>
