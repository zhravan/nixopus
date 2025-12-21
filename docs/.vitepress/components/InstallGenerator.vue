<script setup lang="ts">
import { useInstallGenerator } from './install-generator/useInstallGenerator'
import { CATEGORIES } from './install-generator/constants'
import './install-generator/styles.css'

// Use the composable for all logic
const {
  features,
  copied,
  showAdvanced,
  installMode,
  activeTooltip,
  dismissedInfoMessages,
  draggedId,
  dragOverId,
  hasCustomizations,
  activeCount,
  runCommand,
  validationErrors,
  hasErrors,
  hasWarnings,
  hasInfo,
  isHostIpDisabled,
  areDomainsDisabled,
  isDbPortDisabled,
  isRedisPortDisabled,
  isHealthCheckTimeoutDisabled,
  isAdminEmailDisabled,
  isAdminPasswordDisabled,
  isForceDisabled,
  toggleFeature,
  isFeatureDisabled,
  getFeaturesByCategory,
  copyToClipboard,
  handleDragStart,
  handleDragOver,
  handleDragLeave,
  handleDrop,
  handleDragEnd,
  toggleTooltip,
  closeTooltip,
  dismissInfoMessage,
  resetAll,
} = useInstallGenerator()
</script>

<template>
  <div class="install-generator">

    <!-- Command Output -->
    <div class="command-section">
      <div class="section-header">
        <div class="header-left">
          <h3>Your Install Command</h3>
          <span v-if="activeCount > 0" class="config-count">
            {{ activeCount }} option{{ activeCount > 1 ? 's' : '' }}
          </span>
        </div>
        <button class="reset-btn" @click="resetAll" v-if="hasCustomizations">
          Reset
        </button>
      </div>

      <div class="command-notice">
        <div class="command-content">
          <code>{{ runCommand }}</code>
        </div>
        <button
          class="copy-mini"
          :class="{ copied: copied === 'run' }"
          @click="copyToClipboard(runCommand, 'run')"
          :title="copied === 'run' ? 'Copied!' : 'Copy command'"
        >
          <svg v-if="copied === 'run'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
        </button>
      </div>

      <div class="whats-included">
        <span class="included-label">Installs:</span>
        <div class="included-list">
          <span class="included-item">Nixopus API</span>
          <span class="included-item">Dashboard</span>
          <span class="included-item">PostgreSQL</span>
          <span class="included-item">Redis</span>
          <span class="included-item">Caddy</span>
          <span class="included-item">Docker</span>
        </div>
      </div>

      <!-- Validation Messages -->
      <div v-if="validationErrors.length > 0" class="validation-messages">
        <div
          v-for="(error, index) in validationErrors"
          :key="index"
          class="validation-item"
          :class="error.type"
        >
          <svg v-if="error.type === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <svg v-else-if="error.type === 'warning'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="12" x2="12" y2="16"></line>
            <line x1="12" y1="8" x2="12.01" y2="8"></line>
          </svg>
          <span>{{ error.message }}</span>
          <button
            v-if="error.type === 'info'"
            class="dismiss-btn"
            @click="dismissInfoMessage(error.message)"
            title="Dismiss this message"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Grid -->
    <div class="config-grid">
      
      <!-- Network -->
      <div class="config-category">
        <div class="category-header">
          <h3>{{ CATEGORIES.network.title }}</h3>
          <span class="category-desc">{{ CATEGORIES.network.description }}</span>
        </div>
        <div class="features-list">
          <div 
            v-for="(feature, index) in getFeaturesByCategory('network')" 
            :key="feature.id"
            class="feature-item"
            :class="{ 
              active: feature.enabled,
              disabled: (feature.id === 'hostIp' && isHostIpDisabled) || 
                       ((feature.id === 'apiDomain' || feature.id === 'viewDomain') && areDomainsDisabled),
              dragging: draggedId === feature.id,
              'drag-over': dragOverId === feature.id
            }"
            draggable="true"
            @dragstart="handleDragStart(feature.id)"
            @dragover="handleDragOver($event, feature.id)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, feature.id)"
            @dragend="handleDragEnd"
          >
            <div class="feature-toggle" @click="toggleFeature(feature)">
              <div class="feature-info">
                <div class="feature-label-row">
                  <span class="feature-label">{{ feature.label }}</span>
                  <button
                    v-if="feature.tooltip"
                    class="tooltip-icon"
                    @click.stop="toggleTooltip(feature.id)"
                    :title="feature.tooltip"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="12" cy="12" r="10"></circle>
                      <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                      <line x1="12" y1="17" x2="12.01" y2="17"></line>
                    </svg>
                  </button>
                </div>
                <span class="feature-desc">{{ feature.description }}</span>
                <div v-if="activeTooltip === feature.id && feature.tooltip" class="tooltip-content">
                  {{ feature.tooltip }}
                </div>
              </div>
              <div class="toggle-switch" :class="{ on: feature.enabled }">
                <div class="toggle-knob"></div>
              </div>
            </div>
            <input
              v-if="feature.enabled && feature.inputType !== 'toggle'"
              :type="feature.inputType"
              v-model="feature.value"
              :placeholder="feature.placeholder"
              class="feature-input"
              @click.stop
            />
          </div>
        </div>
      </div>

      <!-- Database -->
      <div class="config-category">
        <div class="category-header">
          <h3>{{ CATEGORIES.database.title }}</h3>
          <span class="category-desc">{{ CATEGORIES.database.description }}</span>
        </div>
        <div class="features-list">
          <div 
            v-for="(feature, index) in getFeaturesByCategory('database')" 
            :key="feature.id"
            class="feature-item"
            :class="{ 
              active: feature.enabled,
              dragging: draggedId === feature.id,
              'drag-over': dragOverId === feature.id
            }"
            draggable="true"
            @dragstart="handleDragStart(feature.id)"
            @dragover="handleDragOver($event, feature.id)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, feature.id)"
            @dragend="handleDragEnd"
          >
            <div class="feature-toggle" @click="toggleFeature(feature)">
              <div class="feature-info">
                <div class="feature-label-row">
                  <span class="feature-label">{{ feature.label }}</span>
                  <button
                    v-if="feature.tooltip"
                    class="tooltip-icon"
                    @click.stop="toggleTooltip(feature.id)"
                    :title="feature.tooltip"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <circle cx="12" cy="12" r="10"></circle>
                      <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                      <line x1="12" y1="17" x2="12.01" y2="17"></line>
                    </svg>
                  </button>
                </div>
                <span class="feature-desc">{{ feature.description }}</span>
                <div v-if="activeTooltip === feature.id && feature.tooltip" class="tooltip-content">
                  {{ feature.tooltip }}
                </div>
              </div>
              <div class="toggle-switch" :class="{ on: feature.enabled }">
                <div class="toggle-knob"></div>
              </div>
            </div>
            <input
              v-if="feature.enabled && feature.inputType !== 'toggle'"
              :type="feature.inputType === 'number' ? 'number' : 'text'"
              v-model="feature.value"
              :placeholder="feature.placeholder"
              class="feature-input"
              @click.stop
            />
          </div>
        </div>
        <div class="category-note">
          Leave disabled to use local Docker database
        </div>
      </div>

      <!-- Options -->
      <div class="config-category options-full-width">
        <div class="category-header">
          <h3>{{ CATEGORIES.options.title }}</h3>
          <span class="category-desc">{{ CATEGORIES.options.description }}</span>
        </div>
        <div class="features-grid">
          <div
            v-for="(feature, index) in getFeaturesByCategory('options')"
            :key="feature.id"
            class="feature-chip"
            :class="{
              active: feature.enabled,
              disabled: isFeatureDisabled(feature),
              dragging: draggedId === feature.id,
              'drag-over': dragOverId === feature.id
            }"
            @click="toggleFeature(feature)"
            draggable="true"
            @dragstart="handleDragStart(feature.id)"
            @dragover="handleDragOver($event, feature.id)"
            @dragleave="handleDragLeave"
            @drop="handleDrop($event, feature.id)"
            @dragend="handleDragEnd"
            :title="feature.tooltip"
          >
            <span class="chip-label">{{ feature.label }}</span>
            <button
              v-if="feature.tooltip"
              class="tooltip-icon-chip"
              @click.stop="toggleTooltip(feature.id)"
            >
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
            </button>
          </div>
        </div>
        <div v-if="getFeaturesByCategory('options').some(f => activeTooltip === f.id)" class="tooltip-popover">
          <div v-for="feature in getFeaturesByCategory('options')" :key="feature.id">
            <div v-if="activeTooltip === feature.id" class="tooltip-content">
              <strong>{{ feature.label }}:</strong> {{ feature.tooltip }}
            </div>
          </div>
        </div>
        <div class="timeout-input" v-if="getFeaturesByCategory('options').find(f => f.id === 'timeout')?.enabled">
          <label>Timeout (seconds)</label>
          <input 
            type="number" 
            v-model="getFeaturesByCategory('options').find(f => f.id === 'timeout')!.value"
            placeholder="300"
          />
        </div>
        <div class="timeout-input" v-if="getFeaturesByCategory('options').find(f => f.id === 'healthCheckTimeout')?.enabled">
          <label>Health Check Timeout (seconds)</label>
          <input 
            type="number" 
            v-model="getFeaturesByCategory('options').find(f => f.id === 'healthCheckTimeout')!.value"
            placeholder="120"
            :disabled="isHealthCheckTimeoutDisabled"
          />
        </div>
        <div v-if="getFeaturesByCategory('options').find(f => f.id === 'configFile')?.enabled" class="feature-item" style="margin-top: 0.75rem;">
          <input
            type="text"
            v-model="getFeaturesByCategory('options').find(f => f.id === 'configFile')!.value"
            placeholder="/path/to/config.yaml"
            class="feature-input"
          />
        </div>
        <div v-if="getFeaturesByCategory('options').find(f => f.id === 'adminEmail')?.enabled" class="feature-item" :class="{ disabled: isAdminEmailDisabled }" style="margin-top: 0.75rem;">
          <input
            type="text"
            v-model="getFeaturesByCategory('options').find(f => f.id === 'adminEmail')!.value"
            placeholder="admin@example.com"
            class="feature-input"
            :disabled="isAdminEmailDisabled"
          />
        </div>
        <div v-if="getFeaturesByCategory('options').find(f => f.id === 'adminPassword')?.enabled" class="feature-item" :class="{ disabled: isAdminPasswordDisabled }" style="margin-top: 0.75rem;">
          <input
            type="password"
            v-model="getFeaturesByCategory('options').find(f => f.id === 'adminPassword')!.value"
            placeholder="secure-password"
            class="feature-input"
            :disabled="isAdminPasswordDisabled"
          />
        </div>
      </div>

    </div>

    <!-- Advanced Toggle -->
    <button class="advanced-toggle" @click="showAdvanced = !showAdvanced">
      <span class="toggle-icon">{{ showAdvanced ? '−' : '+' }}</span>
      {{ showAdvanced ? 'Hide' : 'Show' }} Advanced Options
      <span class="advanced-badge" v-if="!showAdvanced">Ports, Developer</span>
    </button>

    <!-- Advanced Options -->
    <div v-if="showAdvanced" class="advanced-section">
      <div class="config-grid">
        
        <!-- Ports -->
        <div class="config-category compact">
          <div class="category-header">
            <h3>{{ CATEGORIES.ports.title }}</h3>
            <span class="category-desc">{{ CATEGORIES.ports.description }}</span>
          </div>
          <div class="ports-list">
            <div 
              v-for="(feature, index) in getFeaturesByCategory('ports')" 
              :key="feature.id"
              class="port-feature-item"
              :class="{ 
                active: feature.enabled,
                disabled: (feature.id === 'dbPort' && isDbPortDisabled) || 
                         (feature.id === 'redisPort' && isRedisPortDisabled),
                dragging: draggedId === feature.id,
                'drag-over': dragOverId === feature.id
              }"
              draggable="true"
              @dragstart="handleDragStart(feature.id)"
              @dragover="handleDragOver($event, feature.id)"
              @dragleave="handleDragLeave"
              @drop="handleDrop($event, feature.id)"
              @dragend="handleDragEnd"
            >
              <div class="port-toggle" @click="toggleFeature(feature)">
                <div class="port-info">
                  <div class="feature-label-row">
                    <span class="port-label">{{ feature.label }}</span>
                    <button
                      v-if="feature.tooltip"
                      class="tooltip-icon"
                      @click.stop="toggleTooltip(feature.id)"
                      :title="feature.tooltip"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                        <line x1="12" y1="17" x2="12.01" y2="17"></line>
                      </svg>
                    </button>
                  </div>
                  <span class="port-desc">{{ feature.description }}</span>
                  <div v-if="activeTooltip === feature.id && feature.tooltip" class="tooltip-content">
                    {{ feature.tooltip }}
                  </div>
                </div>
                <div class="toggle-switch small" :class="{ on: feature.enabled }">
                  <div class="toggle-knob"></div>
                </div>
              </div>
              <input
                v-if="feature.enabled"
                type="number"
                v-model="feature.value"
                :placeholder="feature.placeholder"
                class="port-input"
                @click.stop
              />
            </div>
          </div>
        </div>

        <!-- Developer -->
        <div class="config-category compact">
          <div class="category-header">
            <h3>{{ CATEGORIES.developer.title }}</h3>
          </div>
          <div class="features-list compact">
            <div 
              v-for="(feature, index) in getFeaturesByCategory('developer')" 
              :key="feature.id"
              class="feature-item compact"
              :class="{ 
                active: feature.enabled,
                dragging: draggedId === feature.id,
                'drag-over': dragOverId === feature.id
              }"
              draggable="true"
              @dragstart="handleDragStart(feature.id)"
              @dragover="handleDragOver($event, feature.id)"
              @dragleave="handleDragLeave"
              @drop="handleDrop($event, feature.id)"
              @dragend="handleDragEnd"
            >
              <div class="feature-toggle" @click="toggleFeature(feature)">
                <div class="feature-info-compact">
                  <div class="feature-label-row">
                    <span class="feature-label">{{ feature.label }}</span>
                    <button
                      v-if="feature.tooltip"
                      class="tooltip-icon"
                      @click.stop="toggleTooltip(feature.id)"
                      :title="feature.tooltip"
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                        <line x1="12" y1="17" x2="12.01" y2="17"></line>
                      </svg>
                    </button>
                  </div>
                  <div v-if="activeTooltip === feature.id && feature.tooltip" class="tooltip-content">
                    {{ feature.tooltip }}
                  </div>
                </div>
                <div class="toggle-switch small" :class="{ on: feature.enabled }">
                  <div class="toggle-knob"></div>
                </div>
              </div>
              <input
                v-if="feature.enabled && feature.inputType !== 'toggle'"
                :type="feature.inputType === 'number' ? 'number' : 'text'"
                v-model="feature.value"
                :placeholder="feature.placeholder"
                class="feature-input"
                @click.stop
              />
            </div>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>
