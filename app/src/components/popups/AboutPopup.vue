<template>
  <div class="flex flex-col min-h-0 flex-grow">
    <h1 class="flex pb-2 mb-2 border-b-1 dark:border-dark-100 space-x-2 items-center">
      <i-ph-info />
      <span class="text-lg">Kiel-Live</span>
    </h1>
    <div class="flex flex-col gap-2 mx-2 text-l overflow-y-auto">
      <p class="mb-2">
        Wo bleibt mein Bus? Kiel-Live zeigt aktuelle Abfahrtszeiten für alle Haltestellen und die aktuelle Position von
        deinem Bus an.
      </p>
      <p class="mb-2">Unterstütze uns indem du die App mit deinen Freunden teilst und gib uns Feedback!</p>
      <Button href="https://www.instagram.com/kiel.live/"><i-mdi-instagram class="mr-2" />Follow @kiel.live</Button>
      <Button href="https://github.com/kiel-live/kiel-live"><i-mdi-github class="mr-2" />Develop on Github</Button>
      <h2 class="border-b-1 mb-2 font-bold">Changelog</h2>
      <h3 class="border-b-1 mt-2 mb-1">{{ new Date('2022-05-03').toLocaleDateString() }}</h3>
      <span>Komplett neue App mit Fokus auf der Karte</span>
      <h3 class="border-b-1 mt-2 mb-1">{{ new Date('2020-05-13').toLocaleDateString() }}</h3>
      <span>Karte mit Haltestellen und Echtzeit-Positionen aller Busse</span>
      <h3 class="border-b-1 mt-2 mb-1">{{ new Date('2019-06-22').toLocaleDateString() }}</h3>
      <span>Anzeige von Verkehrs-Warnungen hinzugefügt</span>
      <h3 class="border-b-1 mt-2 mb-1">{{ new Date('2019-05-20').toLocaleDateString() }}</h3>
      <span><i-fa-solid-birthday-cake class="align-text-bottom" /> Projekt gestartet</span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

import Button from '~/components/atomic/Button.vue';

export default defineComponent({
  name: 'AboutPopup',
  components: { Button },
});
</script>
