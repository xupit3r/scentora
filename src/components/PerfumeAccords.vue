<script setup lang="js">
  const { accords } = defineProps({
    accords: {
      type: Array,
      default: () => []
    }
  });
</script>

<template>
  <div class="accords">
    <h2>Accords</h2>
    <div v-for="accord in accords"
         class="accord">{{ accord }}</div>
  </div>
</template>

<style lang="css">

</style>