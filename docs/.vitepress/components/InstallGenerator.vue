<script setup lang="ts">
import { useInstallGenerator } from './install-generator/composables'
import { CATEGORIES } from './install-generator/utils'
import CommandOutput from './install-generator/sections/CommandOutput.vue'
import AdvancedSection from './install-generator/sections/AdvancedSection.vue'
import OptionsSection from './install-generator/sections/OptionsSection.vue'
import DatabaseSection from './install-generator/sections/DatabaseSection.vue'
import './install-generator/styles.css'

// Use the composable for all logic
const {
  features,
  copied,
  showAdvanced,
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
  dismissInfoMessage,
  resetAll,
} = useInstallGenerator()
</script>

<template>
  <div class="install-generator">

    <!-- Command Output -->
    <CommandOutput
      :run-command="runCommand"
      :active-count="activeCount"
      :has-customizations="hasCustomizations"
      :copied="copied"
      :validation-errors="validationErrors"
      @copy-to-clipboard="copyToClipboard"
      @reset-all="resetAll"
      @dismiss-info-message="dismissInfoMessage"
    />

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
                    @mouseenter="activeTooltip = feature.id"
                    @mouseleave="activeTooltip = null"
                    @click.stop
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
      <DatabaseSection
        :active-tooltip="activeTooltip"
        :dragged-id="draggedId"
        :drag-over-id="dragOverId"
        :get-features-by-category="getFeaturesByCategory"
        :toggle-feature="toggleFeature"
        :handle-drag-start="handleDragStart"
        :handle-drag-over="handleDragOver"
        :handle-drag-leave="handleDragLeave"
        :handle-drop="handleDrop"
        :handle-drag-end="handleDragEnd"
        @update:active-tooltip="activeTooltip = $event"
      />

      <!-- Options -->
      <OptionsSection
        :active-tooltip="activeTooltip"
        :dragged-id="draggedId"
        :drag-over-id="dragOverId"
        :is-feature-disabled="isFeatureDisabled"
        :is-health-check-timeout-disabled="isHealthCheckTimeoutDisabled"
        :is-admin-email-disabled="isAdminEmailDisabled"
        :is-admin-password-disabled="isAdminPasswordDisabled"
        :get-features-by-category="getFeaturesByCategory"
        :toggle-feature="toggleFeature"
        :handle-drag-start="handleDragStart"
        :handle-drag-over="handleDragOver"
        :handle-drag-leave="handleDragLeave"
        :handle-drop="handleDrop"
        :handle-drag-end="handleDragEnd"
        @update:active-tooltip="activeTooltip = $event"
      />

    </div>

    <!-- Advanced Section -->
    <AdvancedSection
      :show-advanced="showAdvanced"
      :active-tooltip="activeTooltip"
      :dragged-id="draggedId"
      :drag-over-id="dragOverId"
      :is-db-port-disabled="isDbPortDisabled"
      :is-redis-port-disabled="isRedisPortDisabled"
      :get-features-by-category="getFeaturesByCategory"
      :toggle-feature="toggleFeature"
      :handle-drag-start="handleDragStart"
      :handle-drag-over="handleDragOver"
      :handle-drag-leave="handleDragLeave"
      :handle-drop="handleDrop"
      :handle-drag-end="handleDragEnd"
      @update:show-advanced="showAdvanced = $event"
      @update:active-tooltip="activeTooltip = $event"
    />

  </div>
</template>
