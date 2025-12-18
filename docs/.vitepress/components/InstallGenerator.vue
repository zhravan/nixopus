<script setup lang="ts">
import { ref, computed, watch } from 'vue'

// Form state - ALL install parameters
const apiDomain = ref('')
const viewDomain = ref('')
const hostIp = ref('')
const repo = ref('')
const branch = ref('')
const configFile = ref('')
const apiPort = ref('')
const viewPort = ref('')
const dbPort = ref('')
const redisPort = ref('')
const caddyAdminPort = ref('')
const caddyHttpPort = ref('')
const caddyHttpsPort = ref('')
const supertokensPort = ref('')
const externalDbUrl = ref('')
const verbose = ref(false)
const dryRun = ref(false)
const force = ref(false)
const staging = ref(false)
const cliOnly = ref(false)
const skipNixopus = ref(false)
const timeout = ref('')

// UI state
const showPorts = ref(false)
const showDeveloper = ref(false)
const copied = ref(false)
const multiLine = ref(false)

// Validation state
const domainError = ref('')
const portError = ref('')

// Computed command
const command = computed(() => {
  const parts: string[] = []

  if (cliOnly.value) {
    if (skipNixopus.value) {
      parts.push('curl -sSL https://install.nixopus.com | bash -s -- --skip-nixopus-install')
    } else {
      parts.push('curl -sSL https://install.nixopus.com | bash')
      
      // Add installer options
      const installerFlags: string[] = []
      
      if (apiDomain.value) installerFlags.push(`--api-domain ${apiDomain.value}`)
      if (viewDomain.value) installerFlags.push(`--view-domain ${viewDomain.value}`)
      if (hostIp.value && !apiDomain.value && !viewDomain.value) {
        installerFlags.push(`--host-ip ${hostIp.value}`)
      }
      if (apiPort.value) installerFlags.push(`--api-port ${apiPort.value}`)
      if (viewPort.value) installerFlags.push(`--view-port ${viewPort.value}`)
      if (dbPort.value && !externalDbUrl.value) installerFlags.push(`--db-port ${dbPort.value}`)
      if (redisPort.value) installerFlags.push(`--redis-port ${redisPort.value}`)
      if (caddyAdminPort.value) installerFlags.push(`--caddy-admin-port ${caddyAdminPort.value}`)
      if (caddyHttpPort.value) installerFlags.push(`--caddy-http-port ${caddyHttpPort.value}`)
      if (caddyHttpsPort.value) installerFlags.push(`--caddy-https-port ${caddyHttpsPort.value}`)
      if (supertokensPort.value) installerFlags.push(`--supertokens-port ${supertokensPort.value}`)
      if (externalDbUrl.value) installerFlags.push(`--external-db-url "${externalDbUrl.value}"`)
      if (repo.value) installerFlags.push(`--repo ${repo.value}`)
      if (branch.value) installerFlags.push(`--branch ${branch.value}`)
      if (staging.value) installerFlags.push('--staging')
      if (verbose.value) installerFlags.push('--verbose')
      if (dryRun.value) installerFlags.push('--dry-run')
      if (force.value) installerFlags.push('--force')
      if (timeout.value) installerFlags.push(`--timeout ${timeout.value}`)
      
      if (installerFlags.length > 0) {
        parts.push('-s --')
        parts.push(...installerFlags)
      }
    }
    return parts.join(multiLine.value ? ' \\\n  ' : ' ')
  }

  parts.push('nixopus install')

  // Configuration file
  if (configFile.value) parts.push(`--config-file ${configFile.value}`)

  // Domains (both required together)
  if (apiDomain.value) parts.push(`--api-domain ${apiDomain.value}`)
  if (viewDomain.value) parts.push(`--view-domain ${viewDomain.value}`)

  // Host IP (mutually exclusive with domains)
  if (hostIp.value && !apiDomain.value && !viewDomain.value) {
    parts.push(`--host-ip ${hostIp.value}`)
  }

  // Ports
  if (apiPort.value) parts.push(`--api-port ${apiPort.value}`)
  if (viewPort.value) parts.push(`--view-port ${viewPort.value}`)
  if (dbPort.value && !externalDbUrl.value) parts.push(`--db-port ${dbPort.value}`)
  if (redisPort.value) parts.push(`--redis-port ${redisPort.value}`)
  if (caddyAdminPort.value) parts.push(`--caddy-admin-port ${caddyAdminPort.value}`)
  if (caddyHttpPort.value) parts.push(`--caddy-http-port ${caddyHttpPort.value}`)
  if (caddyHttpsPort.value) parts.push(`--caddy-https-port ${caddyHttpsPort.value}`)
  if (supertokensPort.value) parts.push(`--supertokens-port ${supertokensPort.value}`)

  // Database
  if (externalDbUrl.value) parts.push(`--external-db-url "${externalDbUrl.value}"`)

  // Repository
  if (repo.value) parts.push(`--repo ${repo.value}`)
  if (branch.value) parts.push(`--branch ${branch.value}`)
  if (staging.value) parts.push('--staging')

  // Options
  if (verbose.value) parts.push('--verbose')
  if (dryRun.value) parts.push('--dry-run')
  if (force.value) parts.push('--force')
  if (timeout.value) parts.push(`--timeout ${timeout.value}`)

  return parts.join(multiLine.value ? ' \\\n  ' : ' ')
})

// Computed display command with syntax highlighting
const displayCommand = computed(() => {
  if (cliOnly.value) {
    const flags: { flag: string; value: string }[] = []
    
    if (skipNixopus.value) {
      return {
        base: 'curl -sSL https://install.nixopus.com | bash -s -- --skip-nixopus-install',
        flags: []
      }
    }
    
    // Build curl command with installer flags
    if (apiDomain.value) flags.push({ flag: '--api-domain', value: apiDomain.value })
    if (viewDomain.value) flags.push({ flag: '--view-domain', value: viewDomain.value })
    if (hostIp.value && !apiDomain.value && !viewDomain.value) {
      flags.push({ flag: '--host-ip', value: hostIp.value })
    }
    if (apiPort.value) flags.push({ flag: '--api-port', value: apiPort.value })
    if (viewPort.value) flags.push({ flag: '--view-port', value: viewPort.value })
    if (dbPort.value && !externalDbUrl.value) flags.push({ flag: '--db-port', value: dbPort.value })
    if (redisPort.value) flags.push({ flag: '--redis-port', value: redisPort.value })
    if (caddyAdminPort.value) flags.push({ flag: '--caddy-admin-port', value: caddyAdminPort.value })
    if (caddyHttpPort.value) flags.push({ flag: '--caddy-http-port', value: caddyHttpPort.value })
    if (caddyHttpsPort.value) flags.push({ flag: '--caddy-https-port', value: caddyHttpsPort.value })
    if (supertokensPort.value) flags.push({ flag: '--supertokens-port', value: supertokensPort.value })
    if (externalDbUrl.value) flags.push({ flag: '--external-db-url', value: `"${externalDbUrl.value}"` })
    if (repo.value) flags.push({ flag: '--repo', value: repo.value })
    if (branch.value) flags.push({ flag: '--branch', value: branch.value })
    if (staging.value) flags.push({ flag: '--staging', value: '' })
    if (verbose.value) flags.push({ flag: '--verbose', value: '' })
    if (dryRun.value) flags.push({ flag: '--dry-run', value: '' })
    if (force.value) flags.push({ flag: '--force', value: '' })
    if (timeout.value) flags.push({ flag: '--timeout', value: timeout.value })
    
    return {
      base: 'curl -sSL https://install.nixopus.com | bash',
      flags: flags.length > 0 ? [{ flag: '-s --', value: '' }, ...flags] : []
    }
  }

  const flags: { flag: string; value: string }[] = []

  if (configFile.value) flags.push({ flag: '--config-file', value: configFile.value })
  if (apiDomain.value) flags.push({ flag: '--api-domain', value: apiDomain.value })
  if (viewDomain.value) flags.push({ flag: '--view-domain', value: viewDomain.value })
  if (hostIp.value && !apiDomain.value && !viewDomain.value) {
    flags.push({ flag: '--host-ip', value: hostIp.value })
  }
  if (apiPort.value) flags.push({ flag: '--api-port', value: apiPort.value })
  if (viewPort.value) flags.push({ flag: '--view-port', value: viewPort.value })
  if (dbPort.value && !externalDbUrl.value) flags.push({ flag: '--db-port', value: dbPort.value })
  if (redisPort.value) flags.push({ flag: '--redis-port', value: redisPort.value })
  if (caddyAdminPort.value) flags.push({ flag: '--caddy-admin-port', value: caddyAdminPort.value })
  if (caddyHttpPort.value) flags.push({ flag: '--caddy-http-port', value: caddyHttpPort.value })
  if (caddyHttpsPort.value) flags.push({ flag: '--caddy-https-port', value: caddyHttpsPort.value })
  if (supertokensPort.value) flags.push({ flag: '--supertokens-port', value: supertokensPort.value })
  if (externalDbUrl.value) flags.push({ flag: '--external-db-url', value: `"${externalDbUrl.value}"` })
  if (repo.value) flags.push({ flag: '--repo', value: repo.value })
  if (branch.value) flags.push({ flag: '--branch', value: branch.value })
  if (staging.value) flags.push({ flag: '--staging', value: '' })
  if (verbose.value) flags.push({ flag: '--verbose', value: '' })
  if (dryRun.value) flags.push({ flag: '--dry-run', value: '' })
  if (force.value) flags.push({ flag: '--force', value: '' })
  if (timeout.value) flags.push({ flag: '--timeout', value: timeout.value })

  return {
    base: 'nixopus install',
    flags
  }
})

// Validation: Domain mismatch (both required together)
watch([apiDomain, viewDomain], ([api, view]) => {
  if ((api && !view) || (!api && view)) {
    domainError.value = 'Both API and View domains required together'
  } else {
    domainError.value = ''
  }
})

// Validation: Port uniqueness
watch([apiPort, viewPort, dbPort, redisPort, caddyAdminPort, caddyHttpPort, caddyHttpsPort, supertokensPort],
  (ports) => {
    const filledPorts = ports.filter(p => p && p !== '').map(p => parseInt(p as string))
    const uniquePorts = new Set(filledPorts)
    if (filledPorts.length !== uniquePorts.size) {
      portError.value = 'All ports must be unique'
    } else {
      portError.value = ''
    }
  }
)

// Copy function
const copyCommand = () => {
  const plainCommand = command.value.replace(/\\\n\s*/g, ' ')
  navigator.clipboard.writeText(plainCommand)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 1200)
}

// Reset function
const resetForm = () => {
  apiDomain.value = ''
  viewDomain.value = ''
  hostIp.value = ''
  repo.value = ''
  branch.value = ''
  configFile.value = ''
  apiPort.value = ''
  viewPort.value = ''
  dbPort.value = ''
  redisPort.value = ''
  caddyAdminPort.value = ''
  caddyHttpPort.value = ''
  caddyHttpsPort.value = ''
  supertokensPort.value = ''
  externalDbUrl.value = ''
  verbose.value = false
  dryRun.value = false
  force.value = false
  staging.value = false
  cliOnly.value = false
  skipNixopus.value = false
  timeout.value = ''
  showPorts.value = false
  showDeveloper.value = false
  domainError.value = ''
  portError.value = ''
}
</script>

<template>
  <div class="install-generator">
    <!-- Command Preview -->
    <div class="command-preview">
      <div class="command-header">
        <h3>Generated Command</h3>
        <div class="command-controls">
          <button
            class="control-btn"
            @click="resetForm"
            title="Reset all"
          >
            ↻
          </button>
          <button
            class="control-btn"
            :class="{ copied: copied }"
            @click="copyCommand"
            title="Copy command"
          >
            <span v-if="copied">✓</span>
            <span v-else>⧉</span>
          </button>
        </div>
      </div>
      <div class="command-display" :class="{ 'multi-line': multiLine }">
        <span class="prompt">$</span>
        <span class="base-command">{{ displayCommand.base }}</span>
        <template v-if="!cliOnly">
          <template v-for="(item, index) in displayCommand.flags" :key="index">
            <span v-if="multiLine" class="line-break"> \</span>
            <span class="flag">{{ item.flag }}</span>
            <span v-if="item.value" class="value">{{ item.value }}</span>
          </template>
        </template>
      </div>
    </div>

    <!-- Validation Errors -->
    <div v-if="domainError || portError" class="error-box">
      <div v-if="domainError">⚠️ {{ domainError }}</div>
      <div v-if="portError">⚠️ {{ portError }}</div>
    </div>

    <!-- Configuration Form -->
    <div class="config-form">
      <!-- CLI Only Mode -->
      <div class="form-section">
        <h3>Installation Mode</h3>
        <label class="checkbox-label">
          <input type="checkbox" v-model="cliOnly" />
          <span>Use curl installer (instead of nixopus CLI)</span>
        </label>
        <label v-if="cliOnly" class="checkbox-label sub-option">
          <input type="checkbox" v-model="skipNixopus" />
          <span>Skip Nixopus installation (CLI only)</span>
        </label>
      </div>

      <template v-if="!cliOnly || (cliOnly && !skipNixopus)">
        <!-- Config File -->
        <div class="form-section">
          <h3>Configuration</h3>
          <input
            type="text"
            v-model="configFile"
            placeholder="--config-file (custom config path)"
          />
        </div>

        <!-- Domains & Host IP -->
        <div class="form-section">
          <h3>Network Configuration</h3>
          <div class="help-text">
            For HTTPS with custom domains (both required together):
          </div>
          <div class="form-row">
            <input
              type="text"
              v-model="apiDomain"
              placeholder="--api-domain api.nixopus.com"
              :disabled="!!hostIp"
              :class="{ 'has-error': domainError }"
            />
            <input
              type="text"
              v-model="viewDomain"
              placeholder="--view-domain nixopus.com"
              :disabled="!!hostIp"
              :class="{ 'has-error': domainError }"
            />
          </div>
          <div class="divider"><span>OR</span></div>
          <div class="help-text">
            For local IP access (auto-detected if not provided):
          </div>
          <input
            type="text"
            v-model="hostIp"
            placeholder="--host-ip 192.168.1.100 or 10.0.0.154"
            :disabled="!!(apiDomain || viewDomain)"
          />
        </div>

        <!-- Database -->
        <div class="form-section">
          <h3>Database</h3>
          <div class="help-text">
            Optional: Use external PostgreSQL (excludes local DB service)
          </div>
          <input
            type="text"
            v-model="externalDbUrl"
            placeholder='--external-db-url "postgresql://user:password@host:port/dbname?sslmode=require"'
          />
        </div>

        <!-- Ports (Collapsible) -->
        <div class="form-section collapsible">
          <div class="section-header" @click="showPorts = !showPorts">
            <h3>Ports - Must Be Unique</h3>
            <span class="expand-icon">{{ showPorts ? '−' : '+' }}</span>
          </div>
          <div v-show="showPorts" class="form-content">
            <div class="help-text">
              Defaults: API=8443 (prod) or 8080 (dev), View=7443 (prod) or 3000 (dev)
            </div>
            <div class="form-row">
              <input
                type="number"
                v-model="apiPort"
                placeholder="--api-port (default: 8443)"
                :class="{ 'has-error': portError }"
              />
              <input
                type="number"
                v-model="viewPort"
                placeholder="--view-port (default: 7443)"
                :class="{ 'has-error': portError }"
              />
            </div>
            <div class="form-row">
              <input
                type="number"
                v-model="dbPort"
                placeholder="--db-port (default: 5432)"
                :disabled="!!externalDbUrl"
                :class="{ 'has-error': portError }"
              />
              <input
                type="number"
                v-model="redisPort"
                placeholder="--redis-port (default: 6379)"
                :class="{ 'has-error': portError }"
              />
            </div>
            <div class="form-row">
              <input
                type="number"
                v-model="caddyAdminPort"
                placeholder="--caddy-admin-port (default: 2019)"
                :class="{ 'has-error': portError }"
              />
              <input
                type="number"
                v-model="caddyHttpPort"
                placeholder="--caddy-http-port (default: 80)"
                :class="{ 'has-error': portError }"
              />
            </div>
            <div class="form-row">
              <input
                type="number"
                v-model="caddyHttpsPort"
                placeholder="--caddy-https-port (default: 443)"
                :class="{ 'has-error': portError }"
              />
              <input
                type="number"
                v-model="supertokensPort"
                placeholder="--supertokens-port (default: 3567)"
                :class="{ 'has-error': portError }"
              />
            </div>
          </div>
        </div>

        <!-- Options -->
        <div class="form-section">
          <h3>Installation Options</h3>
          <div class="checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="verbose" />
              <span>--verbose (show detailed logs)</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="dryRun" />
              <span>--dry-run (preview without changes)</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="force" />
              <span>--force (replace existing files)</span>
            </label>
          </div>
          <input
            type="number"
            v-model="timeout"
            placeholder="--timeout (default: 300 seconds)"
            min="1"
            max="600"
          />
        </div>

        <!-- Developer Options (Collapsible) -->
        <div class="form-section collapsible">
          <div class="section-header" @click="showDeveloper = !showDeveloper">
            <h3>Developer Options</h3>
            <span class="expand-icon">{{ showDeveloper ? '−' : '+' }}</span>
          </div>
          <div v-show="showDeveloper" class="form-content">
            <div class="help-text">
              Defaults: repo=https://github.com/raghavyuva/nixopus, branch=master
            </div>
            <input
              type="text"
              v-model="repo"
              placeholder="--repo https://github.com/user/fork"
            />
            <input
              type="text"
              v-model="branch"
              placeholder="--branch develop"
            />
            <label class="checkbox-label">
              <input type="checkbox" v-model="staging" />
              <span>--staging (use docker-compose-staging.yml)</span>
            </label>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.install-generator {
  max-width: 1000px;
  margin: 2rem auto;
}

/* Command Preview */
.command-preview {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  margin-bottom: 1rem;
  overflow: hidden;
}

.command-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background: var(--vp-c-bg-alt);
  border-bottom: 1px solid var(--vp-c-divider);
}

.command-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}

.command-controls {
  display: flex;
  gap: 0.5rem;
}

.control-btn {
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  padding: 0.4rem 0.8rem;
  cursor: pointer;
  font-size: 1rem;
  color: var(--vp-c-text-2);
  min-width: 36px;
  height: 36px;
}

.control-btn:hover {
  background: var(--vp-c-bg-alt);
  color: var(--vp-c-text-1);
}

.control-btn.active {
  background: var(--vp-c-brand);
  color: white;
  border-color: var(--vp-c-brand);
}

.control-btn.copied {
  background: #22c55e;
  color: white;
  border-color: #22c55e;
}

.command-display {
  padding: 1.5rem;
  font-family: var(--vp-font-family-mono);
  font-size: 0.9rem;
  line-height: 1.8;
  overflow-x: auto;
}

.command-display.multi-line {
  white-space: pre-wrap;
}

.prompt {
  color: var(--vp-c-brand);
  font-weight: bold;
  margin-right: 0.5rem;
  user-select: none;
}

.base-command {
  color: var(--vp-c-text-1);
  font-weight: 600;
}

.line-break {
  color: var(--vp-c-text-3);
}

.flag {
  color: #3b82f6;
  margin-left: 0.5rem;
}

.value {
  color: #f59e0b;
  margin-left: 0.25rem;
}

/* Error Box */
.error-box {
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 6px;
  padding: 1rem;
  margin-bottom: 1rem;
  color: #92400e;
  font-size: 0.9rem;
}

.dark .error-box {
  background: rgba(245, 158, 11, 0.1);
  border-color: #f59e0b;
  color: #fbbf24;
}

.error-box div {
  margin: 0.25rem 0;
}

/* Form */
.config-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-section {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 1.5rem;
}

.form-section h3 {
  margin: 0 0 1rem 0;
  font-size: 1rem;
  font-weight: 600;
}

.help-text {
  font-size: 0.85rem;
  color: var(--vp-c-text-2);
  margin-bottom: 0.75rem;
  line-height: 1.5;
}

.form-section.collapsible {
  padding: 0;
  overflow: hidden;
}

.section-header {
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  user-select: none;
  background: var(--vp-c-bg-alt);
  border-bottom: 1px solid var(--vp-c-divider);
}

.section-header h3 {
  margin: 0;
  flex: 1;
}

.expand-icon {
  font-size: 1rem;
  font-weight: 600;
  color: var(--vp-c-text-3);
  width: 20px;
  text-align: center;
}

.form-content {
  padding: 1.5rem 1.5rem 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-row:last-child {
  margin-bottom: 0;
}

input[type="text"],
input[type="number"] {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  background: var(--vp-c-bg);
  color: var(--vp-c-text-1);
  font-family: var(--vp-font-family-mono);
  font-size: 0.85rem;
  margin-bottom: 0.75rem;
}

input[type="text"]:last-child,
input[type="number"]:last-child {
  margin-bottom: 0;
}

input:focus {
  outline: 2px solid var(--vp-c-brand);
  outline-offset: -2px;
  border-color: transparent;
}

input:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: var(--vp-c-bg-alt);
}

input.has-error {
  border-color: #f59e0b;
}

.checkbox-group {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
  font-family: var(--vp-font-family-mono);
}

.checkbox-label input[type="checkbox"] {
  width: 1.2rem;
  height: 1.2rem;
  cursor: pointer;
  accent-color: var(--vp-c-brand);
  margin: 0;
}

.checkbox-label.sub-option {
  margin-left: 2rem;
  margin-top: 0.5rem;
  color: var(--vp-c-text-2);
}

.divider {
  position: relative;
  text-align: center;
  margin: 1.5rem 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--vp-c-divider);
}

.divider span {
  position: relative;
  background: var(--vp-c-bg-soft);
  padding: 0 1rem;
  font-size: 0.85rem;
  color: var(--vp-c-text-2);
  font-weight: 600;
}

/* Responsive */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .checkbox-group {
    flex-direction: column;
    gap: 0.75rem;
  }

  .command-display {
    font-size: 0.8rem;
  }
}
</style>
