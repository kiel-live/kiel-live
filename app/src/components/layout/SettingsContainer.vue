<template>
  <!-- TODO: remove id settings-container when android apps are updated to set --safe-area-top CSS variable -->
  <div
    id="settings-container"
    class="mt-(--safe-area-top,0px) flex min-h-0 grow flex-col overflow-y-auto p-4 md:mx-auto md:max-w-200 md:px-8 md:py-6"
  >
    <slot />
  </div>
</template>
