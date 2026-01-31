<template>
  <div class="about-page">
    <div class="about-header">
      <h1>🌸 About Scentora</h1>
      <p class="tagline">Your personal accord inventory and perfume formulation companion</p>
    </div>

    <n-card class="about-section" :bordered="false">
      <h2>What is Scentora?</h2>
      <p>
        Scentora is a specialized inventory management system designed for DIY perfumers and fragrance enthusiasts.
        It helps you organize and track your collection of perfume accords, essential oils, and aromatic materials
        with precision and ease.
      </p>
    </n-card>

    <n-card class="about-section" :bordered="false">
      <h2>Key Features</h2>
      <div class="features-grid">
        <div class="feature-item">
          <div class="feature-icon">🏺</div>
          <h3>Accord Inventory</h3>
          <p>Track your accords with detailed information including volume, supplier, and dilution percentage</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🎯</div>
          <h3>Pyramid Positions</h3>
          <p>Organize accords by their position in the fragrance pyramid: top, middle, or base notes</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🏷️</div>
          <h3>Rich Tagging</h3>
          <p>57+ predefined tags for categorization plus custom tags for flexible organization</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🔍</div>
          <h3>Advanced Search</h3>
          <p>Filter by position, volume, supplier, tags, or full-text search across all fields</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">📊</div>
          <h3>Statistics</h3>
          <p>View comprehensive analytics about your collection including pyramid distribution and tag usage</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">⚠️</div>
          <h3>Low Inventory Alerts</h3>
          <p>Get notified when your accords are running low so you can reorder in time</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">📤</div>
          <h3>Backup & Restore</h3>
          <p>Export your entire collection as JSON and import it back for easy data management</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">⌨️</div>
          <h3>Keyboard Shortcuts</h3>
          <p>Power user features including quick actions with keyboard shortcuts (press N for new accord)</p>
        </div>
      </div>
    </n-card>

    <n-card class="about-section" :bordered="false">
      <h2>Technology Stack</h2>
      <div class="tech-stack">
        <div class="tech-category">
          <h3>Frontend</h3>
          <ul>
            <li>Vue.js 3 (Composition API)</li>
            <li>TypeScript</li>
            <li>Pinia (State Management)</li>
            <li>Naive UI (Component Library)</li>
            <li>Tailwind CSS v4</li>
            <li>Vite (Build Tool)</li>
          </ul>
        </div>
        <div class="tech-category">
          <h3>Backend</h3>
          <ul>
            <li>Go (Golang)</li>
            <li>Echo Framework</li>
            <li>PostgreSQL</li>
            <li>JWT Authentication</li>
            <li>RESTful API</li>
          </ul>
        </div>
        <div class="tech-category">
          <h3>Features</h3>
          <ul>
            <li>Secure authentication</li>
            <li>Refresh token rotation</li>
            <li>Rate limiting</li>
            <li>Invitation-only registration</li>
            <li>Complete data isolation</li>
            <li>Responsive design</li>
          </ul>
        </div>
      </div>
    </n-card>

    <n-card class="about-section" :bordered="false">
      <h2>Design Philosophy</h2>
      <p>
        Scentora follows a <strong>Notion-inspired</strong> design philosophy focusing on:
      </p>
      <ul class="philosophy-list">
        <li><strong>Simplicity:</strong> Clean, uncluttered interface that focuses on your content</li>
        <li><strong>Speed:</strong> Fast, responsive interactions with optimistic UI updates</li>
        <li><strong>Flexibility:</strong> Adaptable system that accommodates different workflows</li>
        <li><strong>Accessibility:</strong> Keyboard navigation, ARIA labels, and proper contrast ratios</li>
        <li><strong>Delight:</strong> Smooth animations and helpful feedback at every step</li>
      </ul>
    </n-card>

    <n-card class="about-section version-info" :bordered="false">
      <div class="version-content">
        <div>
          <h3>Version</h3>
          <p>1.0.0 (Phase 8.5)</p>
        </div>
        <div>
          <h3>Status</h3>
          <p>Active Development</p>
        </div>
        <div>
          <h3>License</h3>
          <p>MIT</p>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { NCard } from 'naive-ui';
</script>

<style scoped>
.about-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px;
}

.about-header {
  text-align: center;
  margin-bottom: 48px;
}

.about-header h1 {
  font-size: 40px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 16px 0;
}

.tagline {
  font-size: 18px;
  color: #787774;
  margin: 0;
}

.about-section {
  margin-bottom: 24px;
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
}

.about-section h2 {
  font-size: 24px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 16px 0;
}

.about-section p {
  color: #37352F;
  line-height: 1.7;
  margin: 0 0 16px 0;
}

.about-section p:last-child {
  margin-bottom: 0;
}

/* Features Grid */
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-top: 24px;
}

.feature-item {
  padding: 20px;
  background: #FAFAFA;
  border-radius: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.feature-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
}

.feature-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.feature-item h3 {
  font-size: 16px;
  font-weight: 600;
  color: #37352F;
  margin: 0 0 8px 0;
}

.feature-item p {
  font-size: 14px;
  color: #787774;
  line-height: 1.6;
  margin: 0;
}

/* Tech Stack */
.tech-stack {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 32px;
  margin-top: 24px;
}

.tech-category h3 {
  font-size: 16px;
  font-weight: 600;
  color: #0F766E;
  margin: 0 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #0F766E;
}

.tech-category ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tech-category li {
  padding: 6px 0;
  color: #37352F;
  font-size: 14px;
}

.tech-category li::before {
  content: "→ ";
  color: #0F766E;
  font-weight: bold;
  margin-right: 8px;
}

/* Philosophy List */
.philosophy-list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0 0;
}

.philosophy-list li {
  padding: 12px 16px;
  margin-bottom: 8px;
  background: #FAFAFA;
  border-left: 4px solid #0F766E;
  border-radius: 4px;
  color: #37352F;
  font-size: 14px;
  line-height: 1.6;
}

.philosophy-list strong {
  color: #0F766E;
  font-weight: 600;
}

/* Version Info */
.version-info {
  background: linear-gradient(135deg, #F7F6F3 0%, #FAFAFA 100%);
}

.version-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 24px;
  text-align: center;
}

.version-content h3 {
  font-size: 14px;
  font-weight: 600;
  color: #787774;
  margin: 0 0 8px 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.version-content p {
  font-size: 18px;
  font-weight: 600;
  color: #0F766E;
  margin: 0;
}

@media (max-width: 768px) {
  .about-page {
    padding: 20px;
  }

  .about-header h1 {
    font-size: 28px;
  }

  .tagline {
    font-size: 16px;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .tech-stack {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}
</style>
