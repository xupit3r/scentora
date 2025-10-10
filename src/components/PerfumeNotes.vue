<script setup lang="js">
  const { notes } = defineProps({
    notes: {
      type: Object,
      default: () => ({})
    }
  });



</script>

<template>
  <div class="notes">
    <h3>Notes</h3>
    <div class="notes-top">
      <h5>top</h5>
      <div v-for="note in notes.top" 
           class="note top">{{ note }}</div>
    </div>
    <div class="notes-middle">
      <h5>middle</h5>
      <div v-for="note in notes.middle" 
           class="note">{{ note }}</div>
    </div>
    <div class="notes-base">
      <h5>base</h5>
      <div v-for="note in notes.base" 
           class="note top">{{ note }}</div>
    </div>
  </div>
</template>

<style lang="css">
</style>