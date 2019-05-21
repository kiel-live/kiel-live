<template>
  <div class="home">
    <TextInput @input="searchStop" :debounce="true" :wait="400" placeholder="Haltestelle" autofocus />
    <div v-if="stops" class="stops">
      <div v-for="stop in stops" :key="stop.id" @click="$router.push({ name: 'stop', params: { stop: stop.id } })" class="stop">{{ stop.name }}</div>
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
