<template>
  <div class="statistics">
    <h2>Collection Statistics</h2>

    <div v-if="loading" class="loading">Loading statistics...</div>

    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else-if="stats" class="stats-content">
      <!-- Overview Cards -->
      <div class="overview-grid">
        <div class="stat-card">
          <div class="stat-value">{{ stats.overview.totalPerfumes }}</div>
          <div class="stat-label">Perfumes</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.overview.uniqueDesigners }}</div>
          <div class="stat-label">Designers</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.overview.uniqueNotes }}</div>
          <div class="stat-label">Unique Notes</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.overview.totalJournalEntries }}</div>
          <div class="stat-label">Journal Entries</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.overview.averageRating }}</div>
          <div class="stat-label">Average Rating</div>
        </div>
        <div class="stat-card">
          <div class="stat-value">{{ stats.pyramidStats.avgNotesPerPerfume }}</div>
          <div class="stat-label">Avg Notes/Perfume</div>
        </div>
      </div>

      <!-- Charts Section -->
      <div class="charts-grid">
        <!-- Top Designers -->
        <div class="chart-card">
          <h3>Top Designers</h3>
          <div class="bar-chart">
            <div
              v-for="item in stats.topDesigners.slice(0, 5)"
              :key="item.designer"
              class="bar-item"
            >
              <div class="bar-label">{{ item.designer }}</div>
              <div class="bar-container">
                <div
                  class="bar"
                  :style="{ width: `${(item.count / maxDesignerCount) * 100}%` }"
                >
                  <span class="bar-value">{{ item.count }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Top Notes -->
        <div class="chart-card">
          <h3>Most Used Notes</h3>
          <div class="bar-chart">
            <div
              v-for="item in stats.topNotes.slice(0, 8)"
              :key="item.note"
              class="bar-item"
            >
              <div class="bar-label">{{ item.note }}</div>
              <div class="bar-container">
                <div
                  class="bar"
                  :style="{ width: `${(item.count / maxNoteCount) * 100}%` }"
                >
                  <span class="bar-value">{{ item.count }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Concentration Distribution -->
        <div class="chart-card">
          <h3>Concentration Types</h3>
          <div class="pie-chart-legend">
            <div
              v-for="(count, concentration) in stats.concentrationDistribution"
              :key="concentration"
              class="legend-item"
            >
              <div class="legend-color" :style="{ background: getConcentrationColor(concentration) }"></div>
              <div class="legend-label">{{ concentration }}</div>
              <div class="legend-value">{{ count }}</div>
            </div>
            <div v-if="Object.keys(stats.concentrationDistribution).length === 0" class="empty-state">
              No concentration data
            </div>
          </div>
        </div>

        <!-- Pyramid Distribution -->
        <div class="chart-card">
          <h3>Notes by Pyramid Level</h3>
          <div class="pyramid-viz">
            <div class="pyramid-level top-level">
              <div class="level-label">Top</div>
              <div class="level-value">{{ stats.pyramidStats.topNotes }}</div>
            </div>
            <div class="pyramid-level middle-level">
              <div class="level-label">Middle</div>
              <div class="level-value">{{ stats.pyramidStats.middleNotes }}</div>
            </div>
            <div class="pyramid-level base-level">
              <div class="level-label">Base</div>
              <div class="level-value">{{ stats.pyramidStats.baseNotes }}</div>
            </div>
          </div>
        </div>

        <!-- Year Distribution -->
        <div v-if="stats.yearDistribution.length > 0" class="chart-card full-width">
          <h3>Collection by Year</h3>
          <div class="timeline-chart">
            <div
              v-for="item in stats.yearDistribution"
              :key="item.year"
              class="timeline-item"
            >
              <div class="timeline-year">{{ item.year }}</div>
              <div class="timeline-bar-container">
                <div
                  class="timeline-bar"
                  :style="{ height: `${(item.count / maxYearCount) * 100}%` }"
                >
                  <span class="timeline-value">{{ item.count }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Export Section -->
      <div class="export-section">
        <h3>Export & Import</h3>
        <div class="export-buttons">
          <button class="btn-primary" @click="handleExport">
            📥 Export Collection
          </button>
          <button class="btn-secondary" @click="triggerImport">
            📤 Import Collection
          </button>
          <input
            ref="fileInput"
            type="file"
            accept=".json"
            style="display: none"
            @change="handleImport"
          />
        </div>
        <p class="export-note">
          Export your collection as JSON for backup or sharing. Import will add perfumes to your existing collection.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { statsService, exportService, type CollectionStats } from '@/services/api';

const stats = ref<CollectionStats | null>(null);
const loading = ref(true);
const error = ref<string | null>(null);
const fileInput = ref<HTMLInputElement>();

const maxDesignerCount = computed(() => {
  if (!stats.value) return 1;
  return Math.max(...stats.value.topDesigners.map(d => d.count));
});

const maxNoteCount = computed(() => {
  if (!stats.value) return 1;
  return Math.max(...stats.value.topNotes.map(n => n.count));
});

const maxYearCount = computed(() => {
  if (!stats.value) return 1;
  return Math.max(...stats.value.yearDistribution.map(y => y.count));
});

function getConcentrationColor(concentration: string): string {
  const colors: Record<string, string> = {
    'Parfum': '#8b5cf6',
    'EDP': '#6b4f9e',
    'EDT': '#a78bfa',
    'EDC': '#c4b5fd',
  };
  return colors[concentration] || '#d1d5db';
}

async function loadStats() {
  try {
    loading.value = true;
    stats.value = await statsService.getCollectionStats();
  } catch (err: any) {
    error.value = err.message || 'Failed to load statistics';
  } finally {
    loading.value = false;
  }
}

async function handleExport() {
  try {
    const blob = await exportService.exportCollection();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `scentora-collection-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  } catch (err: any) {
    error.value = err.message || 'Failed to export collection';
  }
}

function triggerImport() {
  fileInput.value?.click();
}

async function handleImport(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  
  if (!file) return;

  try {
    const result = await exportService.importCollection(file);
    alert(`Import complete!\nPerfumes: ${result.perfumesImported}\nJournal Entries: ${result.journalEntriesImported}\nErrors: ${result.errors.length}`);
    
    // Reload stats
    await loadStats();
    
    // Reset file input
    target.value = '';
  } catch (err: any) {
    error.value = err.message || 'Failed to import collection';
  }
}

onMounted(() => {
  loadStats();
});
</script>

<style scoped>
.statistics {
  width: 100%;
}

h2 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

h3 {
  font-size: 1.2rem;
  color: #6b4f9e;
  margin-bottom: 1rem;
}

.loading,
.error {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.error {
  color: #d32f2f;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 2.5rem;
  font-weight: bold;
  color: #6b4f9e;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.9rem;
  color: #666;
  font-weight: 500;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chart-card.full-width {
  grid-column: 1 / -1;
}

.bar-chart {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.bar-item {
  display: grid;
  grid-template-columns: 120px 1fr;
  gap: 0.5rem;
  align-items: center;
}

.bar-label {
  font-size: 0.9rem;
  color: #555;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bar-container {
  background: #f0f0f0;
  border-radius: 4px;
  height: 24px;
  position: relative;
}

.bar {
  background: linear-gradient(90deg, #6b4f9e, #8b5cf6);
  height: 100%;
  border-radius: 4px;
  min-width: 30px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 0.5rem;
  transition: width 0.3s ease;
}

.bar-value {
  color: white;
  font-size: 0.8rem;
  font-weight: 600;
}

.pie-chart-legend {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.legend-color {
  width: 20px;
  height: 20px;
  border-radius: 4px;
}

.legend-label {
  flex: 1;
  color: #555;
  font-size: 0.9rem;
}

.legend-value {
  font-weight: 600;
  color: #333;
}

.pyramid-viz {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.pyramid-level {
  padding: 1rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.top-level {
  background: linear-gradient(135deg, #fff9e6 0%, #fff3cc 100%);
  width: 100%;
}

.middle-level {
  background: linear-gradient(135deg, #ffe6f0 0%, #ffcce0 100%);
  width: 90%;
  margin: 0 auto;
}

.base-level {
  background: linear-gradient(135deg, #e6f0ff 0%, #cce0ff 100%);
  width: 80%;
  margin: 0 auto;
}

.level-label {
  font-weight: 600;
  color: #555;
}

.level-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #6b4f9e;
}

.timeline-chart {
  display: flex;
  gap: 1rem;
  align-items: flex-end;
  height: 200px;
  padding: 1rem 0;
}

.timeline-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.timeline-bar-container {
  width: 100%;
  height: 150px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.timeline-bar {
  width: 60%;
  background: linear-gradient(180deg, #8b5cf6, #6b4f9e);
  border-radius: 4px 4px 0 0;
  min-height: 20px;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 0.25rem;
  transition: height 0.3s ease;
}

.timeline-value {
  color: white;
  font-size: 0.8rem;
  font-weight: 600;
}

.timeline-year {
  font-size: 0.85rem;
  color: #666;
  font-weight: 500;
}

.export-section {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.export-buttons {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #6b4f9e;
  color: white;
}

.btn-primary:hover {
  background: #5a3f8d;
}

.btn-secondary {
  background: #e0e0e0;
  color: #666;
}

.btn-secondary:hover {
  background: #d0d0d0;
}

.export-note {
  color: #666;
  font-size: 0.9rem;
  margin: 0;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 1rem;
  font-style: italic;
}

@media (max-width: 768px) {
  .bar-item {
    grid-template-columns: 80px 1fr;
  }

  .export-buttons {
    flex-direction: column;
  }
}
</style>
