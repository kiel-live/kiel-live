<template>
  <div class="flex w-full justify-around border-t-1 dark:border-gray-700 drop-shadow">
    <router-link
      :to="{ name: 'home' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name !== 'favorites' && $route.name !== 'settings' }"
      aria-label="Karte"
    >
      <i-carbon-map class="w-7 h-7 mb-1" />
      <span class="mt-auto text-sm">Karte</span>
    </router-link>
    <router-link
      :to="{ name: 'favorites' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name === 'favorites' }"
      aria-label="Favoriten"
    >
      <i-ph-star-fill class="w-7 h-7 mb-1" />
      <span class="mt-auto text-sm">Favoriten</span>
    </router-link>
    <router-link
      :to="{ name: 'settings' }"
      class="flex flex-col items-center w-1/3 p-2"
      :class="{ 'text-blue-400': $route.name === 'settings' }"
      aria-label="Über Kiel-Live"
    >
      <i-ic-baseline-settings class="w-7 h-7 mb-1" />
      <span class="mt-auto text-sm">Einstellungen</span>
    </router-link>
  </div>
</template>
