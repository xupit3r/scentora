<script setup lang="js">
  import PerfumeNotes from './PerfumeNotes.vue';
  import PerfumeAccords from './PerfumeAccords.vue';

  const { perfume } = defineProps({
    perfume: {
      type: Object,
      default: () => ({})
    }
  })
</script>

<template>
  <div class="perfume">
    <h3>{{ perfume.name }}</h3>
    <PerfumeNotes :notes="perfume.notes" />
    <PerfumeAccords :accords="perfume.accords" />
  </div>
</template>

<style lang="css">

</style>