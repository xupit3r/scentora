<template>
  <div class="statistics-page">
    <div class="page-header">
      <h1>Collection Statistics</h1>
      <p class="subtitle">Insights into your accord inventory</p>
    </div>

    <div v-if="loading" class="loading-state">
      <n-spin size="large" />
      <p>Loading statistics...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <n-result status="error" title="Failed to Load" :description="error">
        <template #footer>
          <n-button @click="fetchStats">Retry</n-button>
        </template>
      </n-result>
    </div>

    <div v-else-if="stats" class="stats-content">
      <!-- Overview Cards -->
      <div class="overview-grid">
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="Total Accords" :value="stats.overview.totalAccords" />
        </n-card>
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="Total Volume" :value="`${stats.overview.totalVolume.toFixed(1)} ml`" />
        </n-card>
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="Suppliers" :value="stats.overview.totalSuppliers" />
        </n-card>
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="Unique Tags" :value="stats.overview.totalTags" />
        </n-card>
      </div>

      <!-- Pyramid Distribution -->
      <n-card title="Pyramid Position Distribution" class="chart-card" :bordered="false">
        <div class="pyramid-stats">
          <div class="pyramid-row">
            <div class="pyramid-label">Top Notes</div>
            <div class="pyramid-bar">
              <div 
                class="bar-fill top" 
                :style="{ width: `${pyramidPercentage('top')}%` }"
              >
                <span class="bar-value">{{ stats.pyramidStats.topCount }} ({{ stats.pyramidStats.topVolume.toFixed(1) }}ml)</span>
              </div>
            </div>
          </div>
          <div class="pyramid-row">
            <div class="pyramid-label">Middle Notes</div>
            <div class="pyramid-bar">
              <div 
                class="bar-fill middle" 
                :style="{ width: `${pyramidPercentage('middle')}%` }"
              >
                <span class="bar-value">{{ stats.pyramidStats.middleCount }} ({{ stats.pyramidStats.middleVolume.toFixed(1) }}ml)</span>
              </div>
            </div>
          </div>
          <div class="pyramid-row">
            <div class="pyramid-label">Base Notes</div>
            <div class="pyramid-bar">
              <div 
                class="bar-fill base" 
                :style="{ width: `${pyramidPercentage('base')}%` }"
              >
                <span class="bar-value">{{ stats.pyramidStats.baseCount }} ({{ stats.pyramidStats.baseVolume.toFixed(1) }}ml)</span>
              </div>
            </div>
          </div>
        </div>
      </n-card>

      <!-- Two Column Layout -->
      <div class="two-column-grid">
        <!-- Top Tags -->
        <n-card title="Top Tags" class="chart-card" :bordered="false">
          <div v-if="stats.tagStats.length > 0" class="tag-list">
            <div v-for="item in topTags" :key="item.tag" class="tag-item">
              <n-tag :bordered="false">{{ item.tag }}</n-tag>
              <span class="tag-count">{{ item.count }}</span>
            </div>
          </div>
          <n-empty v-else description="No tags yet" />
        </n-card>

        <!-- Top Suppliers -->
        <n-card title="Top Suppliers" class="chart-card" :bordered="false">
          <div v-if="stats.supplierStats.length > 0" class="supplier-list">
            <div v-for="item in topSuppliers" :key="item.supplier" class="supplier-item">
              <span class="supplier-name">{{ item.supplier }}</span>
              <span class="supplier-count">{{ item.count }} accords</span>
            </div>
          </div>
          <n-empty v-else description="No suppliers recorded" />
        </n-card>
      </div>

      <!-- Volume Stats -->
      <n-card title="Volume Statistics" class="chart-card" :bordered="false">
        <div class="volume-stats">
          <div class="volume-item">
            <div class="volume-label">Average Volume</div>
            <div class="volume-value">{{ stats.volumeStats.averageVolume.toFixed(1) }} ml</div>
          </div>
          <div class="volume-item">
            <div class="volume-label">Minimum</div>
            <div class="volume-value">{{ stats.volumeStats.minVolume.toFixed(1) }} ml</div>
          </div>
          <div class="volume-item">
            <div class="volume-label">Maximum</div>
            <div class="volume-value">{{ stats.volumeStats.maxVolume.toFixed(1) }} ml</div>
          </div>
        </div>
      </n-card>

      <!-- Low Inventory Alert -->
      <n-card v-if="stats.lowInventory.length > 0" title="Low Inventory Alert" class="alert-card" :bordered="false">
        <n-alert type="warning" :bordered="false">
          {{ stats.lowInventory.length }} accord(s) running low on volume
        </n-alert>
        <div class="low-inventory-list">
          <div v-for="item in stats.lowInventory" :key="item.accordId" class="low-inventory-item">
            <span class="accord-name">{{ item.name }}</span>
            <span class="accord-volume">{{ item.volumeMl.toFixed(1) }} ml</span>
            <span v-if="item.supplier" class="accord-supplier">{{ item.supplier }}</span>
          </div>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { NCard, NStatistic, NSpin, NResult, NButton, NTag, NEmpty, NAlert } from 'naive-ui';
import api from '../services/api';

interface OverviewStats {
  totalAccords: number;
  totalVolume: number;
  totalSuppliers: number;
  totalTags: number;
}

interface PyramidStats {
  topCount: number;
  middleCount: number;
  baseCount: number;
  topVolume: number;
  middleVolume: number;
  baseVolume: number;
}

interface TagStat {
  tag: string;
  count: number;
}

interface SupplierStat {
  supplier: string;
  count: number;
}

interface VolumeStats {
  averageVolume: number;
  minVolume: number;
  maxVolume: number;
}

interface LowInventoryItem {
  accordId: string;
  name: string;
  volumeMl: number;
  supplier?: string;
}

interface AccordStatistics {
  overview: OverviewStats;
  pyramidStats: PyramidStats;
  tagStats: TagStat[];
  supplierStats: SupplierStat[];
  volumeStats: VolumeStats;
  lowInventory: LowInventoryItem[];
}

const loading = ref(true);
const error = ref('');
const stats = ref<AccordStatistics | null>(null);

const topTags = computed(() => {
  if (!stats.value) return [];
  return stats.value.tagStats
    .sort((a, b) => b.count - a.count)
    .slice(0, 10);
});

const topSuppliers = computed(() => {
  if (!stats.value) return [];
  return stats.value.supplierStats
    .sort((a, b) => b.count - a.count)
    .slice(0, 5);
});

function pyramidPercentage(position: 'top' | 'middle' | 'base'): number {
  if (!stats.value) return 0;
  const total = stats.value.pyramidStats.topCount + 
                stats.value.pyramidStats.middleCount + 
                stats.value.pyramidStats.baseCount;
  if (total === 0) return 0;
  
  const count = position === 'top' ? stats.value.pyramidStats.topCount :
                position === 'middle' ? stats.value.pyramidStats.middleCount :
                stats.value.pyramidStats.baseCount;
  
  return (count / total) * 100;
}

async function fetchStats() {
  loading.value = true;
  error.value = '';
  try {
    const response = await api.get('/api/stats');
    stats.value = response.data;
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'Failed to load statistics';
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchStats();
});
</script>

<style scoped>
.statistics-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 32px;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 8px 0;
}

.subtitle {
  color: #787774;
  font-size: 16px;
  margin: 0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  gap: 16px;
}

.loading-state p {
  color: #787774;
}

.error-state {
  padding: 40px 20px;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}

.chart-card,
.alert-card {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}

.two-column-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

/* Pyramid Stats */
.pyramid-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pyramid-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pyramid-label {
  min-width: 100px;
  color: #37352F;
  font-weight: 500;
}

.pyramid-bar {
  flex: 1;
  height: 40px;
  background: #F7F6F3;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  display: flex;
  align-items: center;
  padding: 0 12px;
  transition: width 0.3s ease;
  border-radius: 8px;
}

.bar-fill.top {
  background: #FEF3C7;
}

.bar-fill.middle {
  background: #E9D5FF;
}

.bar-fill.base {
  background: #F5E6D3;
}

.bar-value {
  color: #37352F;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

/* Tag List */
.tag-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tag-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #F0F0EE;
}

.tag-item:last-child {
  border-bottom: none;
}

.tag-count {
  color: #787774;
  font-size: 14px;
}

/* Supplier List */
.supplier-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.supplier-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #FAFAFA;
  border-radius: 8px;
}

.supplier-name {
  color: #37352F;
  font-weight: 500;
}

.supplier-count {
  color: #787774;
  font-size: 14px;
}

/* Volume Stats */
.volume-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 24px;
}

.volume-item {
  text-align: center;
  padding: 16px;
  background: #FAFAFA;
  border-radius: 8px;
}

.volume-label {
  color: #787774;
  font-size: 14px;
  margin-bottom: 8px;
}

.volume-value {
  color: #37352F;
  font-size: 24px;
  font-weight: 600;
}

/* Low Inventory */
.low-inventory-list {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.low-inventory-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #FEF3C7;
  border-radius: 8px;
  border-left: 4px solid #F59E0B;
}

.accord-name {
  flex: 1;
  color: #37352F;
  font-weight: 500;
}

.accord-volume {
  color: #92400E;
  font-weight: 600;
}

.accord-supplier {
  color: #787774;
  font-size: 14px;
}

@media (max-width: 768px) {
  .statistics-page {
    padding: 20px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .overview-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .two-column-grid {
    grid-template-columns: 1fr;
  }

  .pyramid-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .pyramid-label {
    min-width: auto;
  }

  .pyramid-bar {
    width: 100%;
  }
}
</style>
