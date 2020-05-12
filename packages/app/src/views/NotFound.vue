<template>
  <div class="notFound">
    <h1>OOPS!</h1>
    <h2>Seite nicht gefunden.</h2>
    <router-link :to="{ name: 'home' }">Zurück zur Startseite</router-link>
  </div>
</template>

<script>
export default {
  name: 'notFound',
};
</script>

<style lang="scss" scoped>
  .notFound {
    margin: auto 0;

    h1 {
      font-size: 5rem;
    }

    h2 {
      font-size: 2rem;
      margin-bottom: 2rem;
    }
  }
</style>
