<template>
  <n-card :bordered="true" class="skeleton-card">
    <div class="skeleton-header">
      <div class="skeleton skeleton-title"></div>
      <div class="skeleton skeleton-badge"></div>
    </div>
    
    <div class="skeleton-body">
      <div class="skeleton skeleton-volume"></div>
      <div class="skeleton skeleton-text"></div>
      <div class="skeleton skeleton-text-short"></div>
      
      <div class="skeleton-tags">
        <div class="skeleton skeleton-tag"></div>
        <div class="skeleton skeleton-tag"></div>
        <div class="skeleton skeleton-tag"></div>
      </div>
    </div>
  </n-card>
</template>

<style scoped>
.skeleton-card {
  height: 100%;
}

.skeleton {
  background: linear-gradient(
    90deg,
    #F7F6F3 0%,
    #E9E9E7 50%,
    #F7F6F3 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.skeleton-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.skeleton-title {
  height: 24px;
  flex: 1;
}

.skeleton-badge {
  height: 20px;
  width: 80px;
  border-radius: 12px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-volume {
  height: 32px;
  width: 100px;
}

.skeleton-text {
  height: 16px;
  width: 100%;
}

.skeleton-text-short {
  height: 16px;
  width: 60%;
}

.skeleton-tags {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.skeleton-tag {
  height: 24px;
  width: 60px;
  border-radius: 6px;
}
</style>
