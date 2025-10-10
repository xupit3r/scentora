<script setup lang="js">
  import Perfume from '@/components/Perfume.vue';
  import { usePerfumeStore } from '@/stores/perfumes';
  
  const { perfume } = usePerfumeStore();
</script>

<template>
  <Perfume :perfume="perfume" />
</template>

<style lang="css" scoped>
</style>