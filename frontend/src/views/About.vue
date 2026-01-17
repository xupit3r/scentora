<template>
  <div class="about">
    <h2>About Scentora</h2>
    <p>A personal perfume cataloging application to track your collection and scent experiences.</p>
    
    <h3>Features</h3>
    <ul>
      <li>📝 Catalog perfumes with detailed information</li>
      <li>🏔️ Track perfume pyramids (top, middle, base notes)</li>
      <li>📔 Keep a journal of your scent experiences</li>
      <li>🔍 Search and filter by notes, designer, and more</li>
    </ul>

    <p class="tech-stack">
      Built with Vue.js, Koa.js, TypeScript, and CouchDB
    </p>
  </div>
</template>

<script setup lang="ts">
// About page
</script>

<style scoped>
.about {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

h2 {
  color: #6b4f9e;
  margin-bottom: 1rem;
}

h3 {
  margin-top: 2rem;
  margin-bottom: 1rem;
  color: #333;
}

p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 1rem;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  padding: 0.5rem 0;
  color: #666;
}

.tech-stack {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e0e0e0;
  font-size: 0.9rem;
  color: #999;
}
</style>
