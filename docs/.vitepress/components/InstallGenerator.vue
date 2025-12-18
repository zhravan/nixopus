<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'

// ============================================
// Feature Toggles
// ============================================
interface FeatureOption {
  id: string
  label: string
  description: string
  enabled: boolean
  inputType?: 'toggle' | 'text' | 'number'
  value?: string
  placeholder?: string
  category: 'network' | 'database' | 'ports' | 'options' | 'developer'
  tooltip?: string
}

const features = ref<FeatureOption[]>([
  // Network
  { id: 'apiDomain', label: 'API Domain', description: 'HTTPS endpoint for API', enabled: false, inputType: 'text', value: '', placeholder: 'api.example.com', category: 'network', tooltip: 'Domain for API server. Requires App Domain for auto HTTPS with Let\'s Encrypt.' },
  { id: 'viewDomain', label: 'App Domain', description: 'HTTPS endpoint for dashboard', enabled: false, inputType: 'text', value: '', placeholder: 'app.example.com', category: 'network', tooltip: 'Domain for dashboard. Requires API Domain for auto HTTPS with Let\'s Encrypt.' },
  { id: 'hostIp', label: 'Host IP', description: 'Server IP (auto-detected if empty)', enabled: false, inputType: 'text', value: '', placeholder: '192.168.1.100', category: 'network', tooltip: 'Custom IP for HTTP-only mode. Mutually exclusive with domains.' },

  // Database
  { id: 'externalDb', label: 'External Database', description: 'Use external PostgreSQL', enabled: false, inputType: 'text', value: '', placeholder: 'postgresql://user:pass@host:5432/db', category: 'database', tooltip: 'Use existing PostgreSQL instead of local DB. Ignores DB/Redis port settings.' },

  // Ports
  { id: 'apiPort', label: 'API Port', description: 'Default: 8443', enabled: false, inputType: 'number', value: '', placeholder: '8443', category: 'ports', tooltip: 'Custom API port. Ports < 1024 need root/sudo privileges.' },
  { id: 'viewPort', label: 'Dashboard Port', description: 'Default: 7443', enabled: false, inputType: 'number', value: '', placeholder: '7443', category: 'ports', tooltip: 'Custom dashboard port. Must differ from API port.' },
  { id: 'dbPort', label: 'Database Port', description: 'Default: 5432', enabled: false, inputType: 'number', value: '', placeholder: '5432', category: 'ports', tooltip: 'Local PostgreSQL port. Ignored when using external database.' },
  { id: 'redisPort', label: 'Redis Port', description: 'Default: 6379', enabled: false, inputType: 'number', value: '', placeholder: '6379', category: 'ports', tooltip: 'Local Redis port. Ignored when using external database.' },

  // Options
  { id: 'verbose', label: 'Verbose', description: 'Show detailed logs', enabled: false, inputType: 'toggle', category: 'options', tooltip: 'Show detailed installation logs for debugging.' },
  { id: 'dryRun', label: 'Dry Run', description: 'Preview without changes', enabled: false, inputType: 'toggle', category: 'options', tooltip: 'Preview installation without making actual changes.' },
  { id: 'force', label: 'Force', description: 'Overwrite existing files', enabled: false, inputType: 'toggle', category: 'options', tooltip: 'Overwrite existing config files. Use with caution.' },
  { id: 'timeout', label: 'Timeout', description: 'Default: 300 seconds', enabled: false, inputType: 'number', value: '', placeholder: '300', category: 'options', tooltip: 'Max wait time in seconds. Default: 300s (5 min).' },

  // Developer
  { id: 'repo', label: 'Custom Repo', description: 'Fork repository URL', enabled: false, inputType: 'text', value: '', placeholder: 'https://github.com/user/fork', category: 'developer', tooltip: 'Install from custom fork instead of official repo.' },
  { id: 'branch', label: 'Branch', description: 'Git branch name', enabled: false, inputType: 'text', value: '', placeholder: 'main', category: 'developer', tooltip: 'Install from specific Git branch. Defaults to "main".' },
  { id: 'staging', label: 'Staging', description: 'Use staging compose file', enabled: false, inputType: 'toggle', category: 'developer', tooltip: 'Use staging config for development/testing only.' },
])

// UI State
const copied = ref<'download' | 'run' | null>(null)
const showAdvanced = ref(false)
const installMode = ref<'curl' | 'nixopus'>('curl')
const activeTooltip = ref<string | null>(null)
const dismissedInfoMessages = ref<Set<string>>(new Set())

// Drag and Drop State
const draggedId = ref<string | null>(null)
const dragOverId = ref<string | null>(null)

// Computed disabled states based on selections
const isHostIpDisabled = computed(() => {
  const apiDomain = features.value.find(f => f.id === 'apiDomain')
  const viewDomain = features.value.find(f => f.id === 'viewDomain')
  // Disabled if ANY domain is enabled
  return apiDomain?.enabled || viewDomain?.enabled
})

const areDomainsDisabled = computed(() => {
  const hostIp = features.value.find(f => f.id === 'hostIp')
  // Disabled if Host IP is enabled
  return !!hostIp?.enabled
})

const isDbPortDisabled = computed(() => {
  const externalDb = features.value.find(f => f.id === 'externalDb')
  // Disabled if External DB is enabled
  return !!externalDb?.enabled
})

const isRedisPortDisabled = computed(() => {
  const externalDb = features.value.find(f => f.id === 'externalDb')
  // Disabled if External DB is enabled
  return !!externalDb?.enabled
})

// Check if a specific feature is disabled
const isFeatureDisabled = (feature: FeatureOption): boolean => {
  switch (feature.id) {
    case 'hostIp':
      return isHostIpDisabled.value || false
    case 'apiDomain':
    case 'viewDomain':
      return areDomainsDisabled.value || false
    case 'dbPort':
      return isDbPortDisabled.value || false
    case 'redisPort':
      return isRedisPortDisabled.value || false
    default:
      return false
  }
}

// Validation
interface ValidationError {
  message: string
  type: 'error' | 'warning' | 'info'
}

const validationErrors = computed<ValidationError[]>(() => {
  const errors: ValidationError[] = []
  
  const apiDomain = features.value.find(f => f.id === 'apiDomain')
  const viewDomain = features.value.find(f => f.id === 'viewDomain')
  const hostIp = features.value.find(f => f.id === 'hostIp')
  const externalDb = features.value.find(f => f.id === 'externalDb')
  const apiPort = features.value.find(f => f.id === 'apiPort')
  const viewPort = features.value.find(f => f.id === 'viewPort')
  const dbPort = features.value.find(f => f.id === 'dbPort')
  const redisPort = features.value.find(f => f.id === 'redisPort')
  const verbose = features.value.find(f => f.id === 'verbose')
  const dryRun = features.value.find(f => f.id === 'dryRun')
  const force = features.value.find(f => f.id === 'force')
  const timeout = features.value.find(f => f.id === 'timeout')
  const repo = features.value.find(f => f.id === 'repo')
  const branch = features.value.find(f => f.id === 'branch')
  const staging = features.value.find(f => f.id === 'staging')
  
  const hasApiDomain = apiDomain?.enabled && apiDomain?.value
  const hasViewDomain = viewDomain?.enabled && viewDomain?.value
  const hasHostIp = hostIp?.enabled && hostIp?.value
  const hasExternalDb = externalDb?.enabled && externalDb?.value
  
  // ============================================
  // CRITICAL RULE 1: Domain Mode vs IP Mode (Mutually Exclusive)
  // ============================================
  // If domains provided -> HTTPS mode with Let's Encrypt (Host IP ignored)
  // If Host IP provided -> HTTP mode with IP (Domains ignored)
  // Cannot mix domains and Host IP
  
  if ((hasApiDomain || hasViewDomain) && hasHostIp) {
    errors.push({
      message: 'Domain mode and Host IP mode are mutually exclusive. Use domains for HTTPS or Host IP for HTTP, not both.',
      type: 'error'
    })
  }
  
  // ============================================
  // CRITICAL RULE 2: Both Domains Required Together
  // ============================================
  // For HTTPS with Let's Encrypt, BOTH domains are mandatory
  
  if (hasApiDomain && !hasViewDomain) {
    errors.push({
      message: 'API Domain requires App Domain. Both domains are mandatory for HTTPS setup with Let\'s Encrypt.',
      type: 'error'
    })
  }
  
  if (hasViewDomain && !hasApiDomain) {
    errors.push({
      message: 'App Domain requires API Domain. Both domains are mandatory for HTTPS setup with Let\'s Encrypt.',
      type: 'error'
    })
  }
  
  // ============================================
  // CRITICAL RULE 3: External DB vs Local DB Ports
  // ============================================
  // External DB ignores local DB/Redis ports
  
  if (hasExternalDb) {
    if (dbPort?.enabled && dbPort?.value) {
      errors.push({
        message: 'Database Port is ignored when using External Database. Remove Database Port or disable External Database.',
        type: 'error'
      })
    }
    
    if (redisPort?.enabled && redisPort?.value) {
      errors.push({
        message: 'Redis Port is ignored when using External Database. External DB manages all data storage.',
        type: 'warning'
      })
    }
  }
  
  // ============================================
  // Domain Format Validation
  // ============================================
  
  if (apiDomain?.enabled && apiDomain?.value) {
    const domainRegex = /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/
    if (!domainRegex.test(apiDomain.value)) {
      errors.push({
        message: 'API Domain must be a valid domain (e.g., api.example.com)',
        type: 'error'
      })
    }
    
    // Cannot use localhost or IP in domain mode
    if (apiDomain.value.includes('localhost') || /\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/.test(apiDomain.value)) {
      errors.push({
        message: 'API Domain cannot be localhost or IP address. Use Host IP mode for local/IP-based setup.',
        type: 'error'
      })
    }
  }
  
  if (viewDomain?.enabled && viewDomain?.value) {
    const domainRegex = /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/
    if (!domainRegex.test(viewDomain.value)) {
      errors.push({
        message: 'App Domain must be a valid domain (e.g., app.example.com)',
        type: 'error'
      })
    }
    
    if (viewDomain.value.includes('localhost') || /\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/.test(viewDomain.value)) {
      errors.push({
        message: 'App Domain cannot be localhost or IP address. Use Host IP mode for local/IP-based setup.',
        type: 'error'
      })
    }
  }
  
  // Same domain check
  if (hasApiDomain && hasViewDomain && apiDomain.value === viewDomain.value) {
    errors.push({
      message: 'API Domain and App Domain must be different (e.g., api.example.com and app.example.com)',
      type: 'error'
    })
  }
  
  // ============================================
  // Host IP Validation
  // ============================================
  
  if (hostIp?.enabled && hostIp?.value) {
    const ipRegex = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/
    const match = hostIp.value.match(ipRegex)
    
    if (!match) {
      errors.push({
        message: 'Host IP must be a valid IPv4 address (e.g., 192.168.1.100)',
        type: 'error'
      })
    } else {
      const octets = match.slice(1, 5).map(Number)
      if (octets.some(o => o > 255)) {
        errors.push({
          message: 'Host IP octets must be 0-255',
          type: 'error'
        })
      }
      
      // Warn about private IPs
      const isPrivate = (octets[0] === 10) || 
                       (octets[0] === 172 && octets[1] >= 16 && octets[1] <= 31) ||
                       (octets[0] === 192 && octets[1] === 168)
      
      if (!isPrivate) {
        errors.push({
          message: 'Host IP appears to be public. Ensure firewall rules are configured.',
          type: 'warning'
        })
      }
      
      if (hostIp.value === '127.0.0.1' || hostIp.value.startsWith('127.')) {
        errors.push({
          message: 'Localhost (127.x.x.x) will only be accessible from this machine.',
          type: 'info'
        })
      }
    }
  }
  
  // Auto-detection info
  if (!hasApiDomain && !hasViewDomain && !hasHostIp) {
    errors.push({
      message: 'Using default HTTP mode. Host IP will be auto-detected. Add domains for HTTPS or specify Host IP.',
      type: 'info'
    })
  }
  
  // ============================================
  // External Database Validation
  // ============================================
  
  if (externalDb?.enabled) {
    if (!externalDb?.value) {
      errors.push({
        message: 'External Database URL is required when enabled',
        type: 'error'
      })
    } else {
      if (!externalDb.value.startsWith('postgresql://') && !externalDb.value.startsWith('postgres://')) {
        errors.push({
          message: 'External Database URL must start with postgresql:// or postgres://',
          type: 'error'
        })
      }
      
      // Basic connection string validation
      const dbUrlPattern = /^postgres(ql)?:\/\/([^:]+):([^@]+)@([^:\/]+):?(\d+)?\/(.+)$/
      if (!dbUrlPattern.test(externalDb.value)) {
        errors.push({
          message: 'Invalid database URL format. Expected: postgresql://user:password@host:port/database',
          type: 'error'
        })
      }
    }
  }
  
  // ============================================
  // Port Validation
  // ============================================
  
  const portFields = [
    { field: apiPort, name: 'API Port', default: 8443, id: 'api' },
    { field: viewPort, name: 'Dashboard Port', default: 7443, id: 'view' },
    { field: dbPort, name: 'Database Port', default: 5432, id: 'db' },
    { field: redisPort, name: 'Redis Port', default: 6379, id: 'redis' }
  ]
  
  const enabledPorts: Array<{ name: string; value: number; id: string }> = []
  
  portFields.forEach(({ field, name, default: defaultPort, id }) => {
    if (field?.enabled && field?.value) {
      const port = parseInt(field.value)
      
      if (isNaN(port) || port < 1 || port > 65535) {
        errors.push({
          message: `${name} must be between 1-65535`,
          type: 'error'
        })
      } else {
        enabledPorts.push({ name, value: port, id })
        
        // Privileged port warning
        if (port < 1024 && port !== 443 && port !== 80) {
          errors.push({
            message: `${name} ${port} requires root/sudo privileges. Recommended: ${defaultPort}`,
            type: 'warning'
          })
        }
        
        // Common port conflicts
        const wellKnown: Record<number, string> = {
          22: 'SSH',
          80: 'HTTP',
          443: 'HTTPS', 
          3306: 'MySQL',
          5432: 'PostgreSQL',
          6379: 'Redis',
          27017: 'MongoDB'
        }
        
        if (wellKnown[port] && id !== 'db' && id !== 'redis') {
          errors.push({
            message: `${name} ${port} is commonly used by ${wellKnown[port]}. May cause conflicts.`,
            type: 'warning'
          })
        }
      }
    }
  })
  
  // Port conflict detection
  const portMap = new Map<number, string[]>()
  enabledPorts.forEach(p => {
    if (!portMap.has(p.value)) {
      portMap.set(p.value, [])
    }
    portMap.get(p.value)!.push(p.name)
  })
  
  portMap.forEach((services, port) => {
    if (services.length > 1) {
      errors.push({
        message: `Port ${port} conflict: ${services.join(' and ')} cannot use the same port`,
        type: 'error'
      })
    }
  })
  
  // Domain mode + custom ports warning
  if (hasApiDomain && hasViewDomain) {
    if (apiPort?.enabled && apiPort?.value) {
      const port = parseInt(apiPort.value)
      if (port !== 443 && port !== 8443) {
        errors.push({
          message: `API Port ${port} with HTTPS domains. Ensure DNS and firewall allow this port.`,
          type: 'info'
        })
      }
    }
    
    if (viewPort?.enabled && viewPort?.value) {
      const port = parseInt(viewPort.value)
      if (port !== 443 && port !== 7443) {
        errors.push({
          message: `Dashboard Port ${port} with HTTPS domains. Ensure DNS and firewall allow this port.`,
          type: 'info'
        })
      }
    }
  }
  
  // ============================================
  // Installation Options Validation
  // ============================================
  
  // Dry Run + Force contradiction
  if (dryRun?.enabled && force?.enabled) {
    errors.push({
      message: 'Dry Run previews changes without executing. Force flag will be ignored in dry run mode.',
      type: 'warning'
    })
  }
  
  // Timeout validation
  if (timeout?.enabled && timeout?.value) {
    const timeoutVal = parseInt(timeout.value)
    
    if (isNaN(timeoutVal) || timeoutVal < 1) {
      errors.push({
        message: 'Timeout must be a positive number (seconds)',
        type: 'error'
      })
    } else if (timeoutVal < 60) {
      errors.push({
        message: 'Timeout < 60s may cause installation failure. Recommended: 300s (5 minutes)',
        type: 'warning'
      })
    } else if (timeoutVal > 3600) {
      errors.push({
        message: 'Timeout > 3600s (1 hour) is unusually long',
        type: 'warning'
      })
    }
  }
  
  // ============================================
  // Developer Options Validation
  // ============================================
  
  // Custom repository
  if (repo?.enabled) {
    if (!repo?.value) {
      errors.push({
        message: 'Repository URL required when enabled',
        type: 'error'
      })
    } else {
      if (!repo.value.startsWith('http://') && !repo.value.startsWith('https://')) {
        errors.push({
          message: 'Repository URL must start with http:// or https://',
          type: 'error'
        })
      }
      
      if (!repo.value.includes('github.com') && !repo.value.includes('gitlab.com')) {
        errors.push({
          message: 'Repository should be from GitHub or GitLab',
          type: 'warning'
        })
      }
      
      // Branch recommendation
      if (!branch?.enabled || !branch?.value) {
        errors.push({
          message: 'Branch name recommended with custom repository. Defaults to "main" if not specified.',
          type: 'info'
        })
      }
    }
  }
  
  // Branch without custom repo
  if (branch?.enabled && branch?.value && (!repo?.enabled || !repo?.value)) {
    errors.push({
      message: 'Branch specified without custom repository. Will use official Nixopus repo with this branch.',
      type: 'info'
    })
  }
  
  // Staging mode validation
  if (staging?.enabled) {
    if (repo?.enabled && repo?.value) {
      errors.push({
        message: 'Staging mode with custom repository. Ensure your fork has docker-compose-staging.yml file.',
        type: 'warning'
      })
    }
    
    if (hasApiDomain && hasViewDomain) {
      errors.push({
        message: 'Staging mode with production domains. Ensure this is intentional.',
        type: 'warning'
      })
    }
  }
  
  // ============================================
  // Empty/Incomplete Field Warnings
  // ============================================
  
  if (apiDomain?.enabled && !apiDomain?.value) {
    errors.push({
      message: 'API Domain is enabled but empty',
      type: 'warning'
    })
  }
  
  if (viewDomain?.enabled && !viewDomain?.value) {
    errors.push({
      message: 'App Domain is enabled but empty',
      type: 'warning'
    })
  }
  
  if (hostIp?.enabled && !hostIp?.value) {
    errors.push({
      message: 'Host IP is enabled but empty. Will be auto-detected.',
      type: 'info'
    })
  }
  
  if (externalDb?.enabled && !externalDb?.value) {
    errors.push({
      message: 'External Database is enabled but URL is empty',
      type: 'warning'
    })
  }
  
  if (repo?.enabled && !repo?.value) {
    errors.push({
      message: 'Custom Repository is enabled but URL is empty',
      type: 'warning'
    })
  }
  
  if (branch?.enabled && !branch?.value) {
    errors.push({
      message: 'Branch is enabled but name is empty',
      type: 'warning'
    })
  }
  
  // ============================================
  // Helpful Info Messages
  // ============================================
  
  if (!hasCustomizations.value) {
    errors.push({
      message: 'Default installation: HTTP mode with auto-detected IP, local PostgreSQL and Redis. Toggle options to customize.',
      type: 'info'
    })
  }
  
  return errors.filter(error => !dismissedInfoMessages.value.has(error.message))
})

const hasErrors = computed(() =>
  validationErrors.value.some(e => e.type === 'error')
)

const hasWarnings = computed(() => 
  validationErrors.value.some(e => e.type === 'warning')
)

const hasInfo = computed(() => 
  validationErrors.value.some(e => e.type === 'info')
)

// Category helpers
const categories = {
  network: { title: 'Network', description: 'Domain & IP configuration' },
  database: { title: 'Database', description: 'PostgreSQL configuration' },
  ports: { title: 'Ports', description: 'Service ports (all optional)' },
  options: { title: 'Options', description: 'Installation behavior' },
  developer: { title: 'Developer', description: 'For contributors & forks' },
}

const getFeaturesByCategory = (category: string) => 
  features.value.filter(f => f.category === category)

// Toggle feature
const toggleFeature = (feature: FeatureOption) => {
  // Check if feature is disabled - prevent enabling
  if (isFeatureDisabled(feature)) return
  
  // If trying to enable, check for conflicts first
  if (!feature.enabled) {
    // Prevent enabling Host IP if any domain is enabled
    if (feature.id === 'hostIp') {
      const apiDomain = features.value.find(f => f.id === 'apiDomain')
      const viewDomain = features.value.find(f => f.id === 'viewDomain')
      if (apiDomain?.enabled || viewDomain?.enabled) return
    }
    
    // Prevent enabling domains if Host IP is enabled
    if (feature.id === 'apiDomain' || feature.id === 'viewDomain') {
      const hostIp = features.value.find(f => f.id === 'hostIp')
      if (hostIp?.enabled) return
    }
    
    // Prevent enabling DB/Redis ports if External DB is enabled
    if (feature.id === 'dbPort' || feature.id === 'redisPort') {
      const externalDb = features.value.find(f => f.id === 'externalDb')
      if (externalDb?.enabled) return
    }
  }
  
  feature.enabled = !feature.enabled
  if (!feature.enabled) {
    feature.value = ''
  }
}

// Commands
const curlCommand = 'curl -sSL https://install.nixopus.com | bash'
const nixopusCommand = 'nixopus install'

const baseCommand = computed(() => 
  installMode.value === 'curl' ? curlCommand : nixopusCommand
)

const runCommand = computed(() => {
  const flags: string[] = []
  
  features.value.forEach(f => {
    if (!f.enabled) return
    
    switch (f.id) {
      case 'apiDomain':
        if (f.value) flags.push(`--api-domain ${f.value}`)
        break
      case 'viewDomain':
        if (f.value) flags.push(`--view-domain ${f.value}`)
        break
      case 'hostIp':
        if (f.value) flags.push(`--host-ip ${f.value}`)
        break
      case 'externalDb':
        if (f.value) flags.push(`--external-db-url "${f.value}"`)
        break
      case 'apiPort':
        if (f.value) flags.push(`--api-port ${f.value}`)
        break
      case 'viewPort':
        if (f.value) flags.push(`--view-port ${f.value}`)
        break
      case 'dbPort':
        if (f.value) flags.push(`--db-port ${f.value}`)
        break
      case 'redisPort':
        if (f.value) flags.push(`--redis-port ${f.value}`)
        break
      case 'verbose':
        flags.push('--verbose')
        break
      case 'dryRun':
        flags.push('--dry-run')
        break
      case 'force':
        flags.push('--force')
        break
      case 'timeout':
        if (f.value) flags.push(`--timeout ${f.value}`)
        break
      case 'repo':
        if (f.value) flags.push(`--repo ${f.value}`)
        break
      case 'branch':
        if (f.value) flags.push(`--branch ${f.value}`)
        break
      case 'staging':
        flags.push('--staging')
        break
    }
  })
  
  if (flags.length === 0) {
    return baseCommand.value
  }
  
  if (installMode.value === 'curl') {
    return `${curlCommand} -s -- ${flags.join(' ')}`
  } else {
    return `${nixopusCommand} ${flags.join(' ')}`
  }
})

const hasCustomizations = computed(() => 
  features.value.some(f => f.enabled)
)

const activeCount = computed(() => 
  features.value.filter(f => f.enabled).length
)

// Copy handlers
const copyToClipboard = async (text: string, type: 'download' | 'run') => {
  try {
    await navigator.clipboard.writeText(text)
    copied.value = type
    setTimeout(() => { copied.value = null }, 2000)
  } catch {
    // Fallback
    const el = document.createElement('textarea')
    el.value = text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
    copied.value = type
    setTimeout(() => { copied.value = null }, 2000)
  }
}

// Drag and Drop Handlers
const handleDragStart = (featureId: string) => {
  draggedId.value = featureId
}

const handleDragOver = (e: DragEvent, featureId: string) => {
  e.preventDefault()
  dragOverId.value = featureId
}

const handleDragLeave = () => {
  dragOverId.value = null
}

const handleDrop = (e: DragEvent, dropFeatureId: string) => {
  e.preventDefault()
  
  if (draggedId.value === null || draggedId.value === dropFeatureId) {
    draggedId.value = null
    dragOverId.value = null
    return
  }
  
  const draggedIndex = features.value.findIndex(f => f.id === draggedId.value)
  const dropIndex = features.value.findIndex(f => f.id === dropFeatureId)
  
  if (draggedIndex === -1 || dropIndex === -1) {
    draggedId.value = null
    dragOverId.value = null
    return
  }
  
  const items = [...features.value]
  const [draggedItem] = items.splice(draggedIndex, 1)
  items.splice(dropIndex, 0, draggedItem)
  
  features.value = items
  draggedId.value = null
  dragOverId.value = null
}

const handleDragEnd = () => {
  draggedId.value = null
  dragOverId.value = null
}

// Tooltip handlers
const toggleTooltip = (featureId: string) => {
  activeTooltip.value = activeTooltip.value === featureId ? null : featureId
}

const closeTooltip = () => {
  activeTooltip.value = null
}

// Dismiss info message
const dismissInfoMessage = (message: string) => {
  dismissedInfoMessages.value.add(message)
}

// Reset
const resetAll = () => {
  features.value.forEach(f => {
    f.enabled = false
    f.value = ''
  })
}
</script>

<template>
  <div class="install-generator">

    <!-- Installation Mode Toggle -->
    <div class="install-mode-selector">
      <div class="mode-label">Installation method:</div>
      <div class="mode-buttons">
        <button 
          class="mode-btn" 
          :class="{ active: installMode === 'curl' }"
          @click="installMode = 'curl'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="16 18 22 12 16 6"></polyline>
            <polyline points="8 6 2 12 8 18"></polyline>
          </svg>
          <span>cURL Script</span>
        </button>
        <button 
          class="mode-btn" 
          :class="{ active: installMode === 'nixopus' }"
          @click="installMode = 'nixopus'"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <line x1="9" y1="9" x2="15" y2="9"></line>
            <line x1="9" y1="15" x2="15" y2="15"></line>
          </svg>
          <span>Nixopus CLI</span>
        </button>
      </div>
    </div>

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
          <h3>{{ categories.network.title }}</h3>
          <span class="category-desc">{{ categories.network.description }}</span>
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
          <h3>{{ categories.database.title }}</h3>
          <span class="category-desc">{{ categories.database.description }}</span>
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
      <div class="config-category">
        <div class="category-header">
          <h3>{{ categories.options.title }}</h3>
          <span class="category-desc">{{ categories.options.description }}</span>
        </div>
        <div class="features-grid">
          <div
            v-for="(feature, index) in getFeaturesByCategory('options')"
            :key="feature.id"
            class="feature-chip"
            :class="{
              active: feature.enabled,
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
            <h3>{{ categories.ports.title }}</h3>
            <span class="category-desc">{{ categories.ports.description }}</span>
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
            <h3>{{ categories.developer.title }}</h3>
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

    <!-- Help -->
    <details class="help-section">
      <summary>Examples & Help</summary>
      <div class="examples">
        <div class="example">
          <h5>Quick Install</h5>
          <code>curl -sSL https://install.nixopus.com | bash</code>
          <p>Auto-detects IP, uses all defaults</p>
        </div>
        <div class="example">
          <h5>Production with HTTPS</h5>
          <code>curl -sSL https://install.nixopus.com | bash -s -- --api-domain api.example.com --view-domain app.example.com</code>
          <p>Requires both domains, auto SSL via Let's Encrypt</p>
        </div>
        <div class="example">
          <h5>Preview Only</h5>
          <code>curl -sSL https://install.nixopus.com | bash -s -- --dry-run --verbose</code>
          <p>See what would happen without making changes</p>
        </div>
      </div>
    </details>

  </div>
</template>

<style scoped>
.install-generator {
  max-width: 800px;
  margin: 0 auto;
}

/* Header */
.generator-header {
  margin-bottom: 1.5rem;
}

.generator-header h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.75rem;
  font-weight: 600;
  line-height: 1.2;
}

.generator-header p {
  margin: 0;
  color: var(--vp-c-text-2);
  font-size: 0.9375rem;
  line-height: 1.5;
}

/* Installation Mode Selector */
.install-mode-selector {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  margin-bottom: 1.5rem;
}

.mode-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  flex-shrink: 0;
}

.mode-buttons {
  display: flex;
  gap: 0.5rem;
  flex: 1;
}

.mode-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  transition: all 0.2s ease;
  flex: 1;
  justify-content: center;
}

.mode-btn:hover {
  border-color: var(--vp-c-brand);
  background: var(--vp-c-bg-soft);
}

.mode-btn.active {
  background: var(--vp-c-default-soft);
  border-color: var(--vp-c-default-2);
  color: var(--vp-c-text-1);
  font-weight: 600;
}

.mode-btn svg {
  flex-shrink: 0;
}

/* Config Grid */
.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

/* Category */
.config-category {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  padding: 1.25rem;
  transition: border-color 0.2s ease;
}

.config-category:hover {
  border-color: var(--vp-c-brand-light);
}

.config-category.compact {
  padding: 1rem;
}

.category-header {
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--vp-c-divider);
}

.category-header h3 {
  margin: 0 0 0.25rem 0;
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.4;
}

.category-desc {
  display: block;
  margin: 0;
  font-size: 0.75rem;
  color: var(--vp-c-text-3);
  line-height: 1.4;
}

.category-note {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--vp-c-divider);
  font-size: 0.75rem;
  color: var(--vp-c-text-3);
  line-height: 1.4;
}

/* Feature List */
.features-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.features-list.compact {
  gap: 0.5rem;
}

.feature-item {
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  overflow: hidden;
}

.feature-item.active {
  border-color: var(--vp-c-brand);
}

.feature-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--vp-c-bg-soft);
}

.feature-item.disabled .feature-toggle {
  cursor: not-allowed;
  pointer-events: none;
}

.feature-toggle {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  cursor: pointer;
}

.feature-info {
  flex: 1;
  min-width: 0;
}

.feature-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.4;
  color: var(--vp-c-text-1);
}

.feature-desc {
  display: block;
  font-size: 0.75rem;
  color: var(--vp-c-text-3);
  line-height: 1.4;
  margin-top: 0.125rem;
}

.feature-input {
  width: 100%;
  padding: 0.625rem 0.75rem;
  border: none;
  border-top: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
  font-size: 0.875rem;
  font-family: var(--vp-font-family-mono);
  color: var(--vp-c-text-1);
  line-height: 1.4;
}

.feature-input:focus {
  outline: none;
  background: var(--vp-c-bg);
}

.feature-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--vp-c-bg-soft);
}

.feature-input::placeholder {
  color: var(--vp-c-text-3);
}

/* Toggle Switch */
.toggle-switch {
  width: 40px;
  height: 20px;
  background: var(--vp-c-divider);
  border-radius: 10px;
  position: relative;
  flex-shrink: 0;
}

.toggle-switch.small {
  width: 32px;
  height: 16px;
  border-radius: 8px;
}

.toggle-switch.on {
  background: var(--vp-c-default-2);
}

.toggle-knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.toggle-switch.small .toggle-knob {
  width: 12px;
  height: 12px;
}

.toggle-switch.on .toggle-knob {
  transform: translateX(20px);
}

.toggle-switch.small.on .toggle-knob {
  transform: translateX(16px);
}

/* Feature Chips */
.features-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.feature-chip {
  padding: 0.5rem 0.875rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  line-height: 1.4;
}

.feature-chip:hover {
  border-color: var(--vp-c-brand);
}

.feature-chip.active {
  background: var(--vp-c-default-2);
  border-color: var(--vp-c-default-2);
  color: white;
}

.timeout-input {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--vp-c-divider);
}

.timeout-input label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
}

.timeout-input input {
  width: 80px;
  padding: 0.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: var(--vp-font-family-mono);
  text-align: center;
  background: var(--vp-c-bg);
}

.timeout-input input:focus {
  outline: none;
  border-color: var(--vp-c-brand);
}

/* Ports Section - Improved Design */
.ports-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.port-feature-item {
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  overflow: hidden;
  transition: border-color 0.2s ease;
}

.port-feature-item.active {
  border-color: var(--vp-c-default-2);
}

.port-feature-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--vp-c-bg-soft);
}

.port-feature-item.disabled .port-toggle {
  cursor: not-allowed;
  pointer-events: none;
}

.port-toggle {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  cursor: pointer;
}

.port-info {
  flex: 1;
  min-width: 0;
}

.port-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1.4;
  color: var(--vp-c-text-1);
}

.port-desc {
  display: block;
  font-size: 0.75rem;
  color: var(--vp-c-text-3);
  line-height: 1.4;
  margin-top: 0.125rem;
}

.port-input {
  width: 100%;
  padding: 0.625rem 0.75rem;
  border: none;
  border-top: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
  font-size: 0.875rem;
  font-family: var(--vp-font-family-mono);
  color: var(--vp-c-text-1);
  text-align: center;
  line-height: 1.4;
}

.port-input:focus {
  outline: none;
  background: var(--vp-c-bg);
}

.port-input::placeholder {
  color: var(--vp-c-text-3);
}

/* Advanced Toggle */
.advanced-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  color: var(--vp-c-text-2);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  margin-bottom: 1rem;
  line-height: 1.4;
}

.advanced-toggle:hover {
  border-color: var(--vp-c-brand);
}

.toggle-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  line-height: 1;
}

.advanced-badge {
  font-size: 0.75rem;
  font-weight: 400;
  color: var(--vp-c-text-3);
}

.advanced-section {
  margin-bottom: 1rem;
}

/* Command Section */
.command-section {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  padding: 1.25rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.copy-mini {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  padding: 0;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  cursor: pointer;
  color: var(--vp-c-text-1);
  flex-shrink: 0;
}

.copy-mini:hover {
  border-color: var(--vp-c-brand);
  background: var(--vp-c-bg-soft);
}

.copy-mini.copied {
  background: #22c55e;
  border-color: #22c55e;
  color: white;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.section-header h3 {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.4;
}

.config-count {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
  background: var(--vp-c-bg);
  padding: 0.375rem 0.625rem;
  border-radius: 4px;
  border: 1px solid var(--vp-c-divider);
  line-height: 1;
}

.reset-btn {
  padding: 0.5rem 1rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
  color: var(--vp-c-text-2);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  line-height: 1;
}

.reset-btn:hover {
  border-color: var(--vp-c-brand);
}

.command-notice {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
  margin-bottom: 1rem;
}

.command-content {
  flex: 1;
  min-width: 0;
}

.command-content code {
  display: block;
  font-family: var(--vp-font-family-mono);
  font-size: 0.8125rem;
  line-height: 1.6;
  word-break: break-all;
  white-space: pre-wrap;
  color: var(--vp-c-text-1);
}

/* What's Included */
.whats-included {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.included-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--vp-c-text-3);
  flex-shrink: 0;
}

.included-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.included-item {
  font-size: 0.75rem;
  color: var(--vp-c-text-2);
  padding: 0.25rem 0.5rem;
  background: var(--vp-c-bg);
  border-radius: 3px;
  line-height: 1.4;
}

/* Validation Messages */
.validation-messages {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--vp-c-divider);
}

.validation-item {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 0.8125rem;
  line-height: 1.5;
}

.validation-item.error {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.validation-item.warning {
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  color: #f59e0b;
}

.validation-item.info {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #3b82f6;
}

.validation-item svg {
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.validation-item span {
  flex: 1;
}

.dismiss-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  padding: 0;
  margin-left: 0.5rem;
  background: transparent;
  border: none;
  cursor: pointer;
  color: currentColor;
  opacity: 0.6;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
}

.dismiss-btn:hover {
  opacity: 1;
}

.dismiss-btn svg {
  display: block;
}

/* Help Section */
.help-section {
  margin-top: 2rem;
  padding: 1rem;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 6px;
}

.help-section summary {
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  color: var(--vp-c-text-2);
  outline: none;
  line-height: 1.4;
}

.help-section summary:hover {
  color: var(--vp-c-text-1);
}

.help-section[open] summary {
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--vp-c-divider);
}

.examples {
  display: grid;
  gap: 1rem;
}

.example {
  padding: 1rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-divider);
  border-radius: 4px;
}

.example h5 {
  margin: 0 0 0.5rem 0;
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  line-height: 1.4;
}

.example code {
  display: block;
  padding: 0.625rem;
  background: var(--vp-c-bg-soft);
  border-radius: 4px;
  font-size: 0.75rem;
  font-family: var(--vp-font-family-mono);
  overflow-x: auto;
  color: var(--vp-c-text-1);
  line-height: 1.5;
}

.example p {
  margin: 0.5rem 0 0 0;
  font-size: 0.75rem;
  color: var(--vp-c-text-3);
  line-height: 1.5;
}

/* Responsive */
@media (max-width: 640px) {
  .install-mode-selector {
    flex-direction: column;
    align-items: stretch;
  }
  
  .mode-buttons {
    flex-direction: column;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
  }
  
  .command-notice {
    flex-direction: column;
    align-items: stretch;
  }
  
  .command-content {
    flex-direction: column;
    align-items: stretch;
  }
  
  .command-content code {
    white-space: normal;
    word-break: break-all;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
  
  .whats-included {
    flex-direction: column;
    align-items: flex-start;
  }
}

/* Drag and Drop Styles */
.feature-item,
.feature-chip,
.port-feature-item {
  cursor: grab;
  transition: all 0.2s ease, opacity 0.15s ease, transform 0.15s ease;
}

.feature-item:active,
.feature-chip:active,
.port-feature-item:active {
  cursor: grabbing;
}

.feature-item.dragging,
.feature-chip.dragging,
.port-feature-item.dragging {
  opacity: 0.4;
  transform: scale(0.95);
}

.feature-item.drag-over,
.feature-chip.drag-over,
.port-feature-item.drag-over {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-color: var(--vp-c-brand);
}

.feature-item.disabled,
.port-feature-item.disabled {
  cursor: not-allowed !important;
}

/* Tooltip Styles */
.feature-label-row {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.tooltip-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  padding: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--vp-c-text-3);
  transition: color 0.2s ease;
  flex-shrink: 0;
}

.tooltip-icon:hover {
  color: var(--vp-c-brand);
}

.tooltip-icon svg {
  display: block;
}

.tooltip-icon-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  padding: 0;
  background: transparent;
  border: none;
  cursor: pointer;
  color: currentColor;
  opacity: 0.7;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
}

.tooltip-icon-chip:hover {
  opacity: 1;
}

.tooltip-icon-chip svg {
  display: block;
}

.tooltip-content {
  margin-top: 0.5rem;
  padding: 0.625rem 0.75rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-brand-light);
  border-radius: 4px;
  font-size: 0.75rem;
  line-height: 1.5;
  color: var(--vp-c-text-2);
  animation: tooltipFadeIn 0.2s ease;
}

.tooltip-popover {
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: var(--vp-c-bg);
  border: 1px solid var(--vp-c-brand-light);
  border-radius: 4px;
  animation: tooltipFadeIn 0.2s ease;
}

.tooltip-popover .tooltip-content {
  margin: 0;
  padding: 0;
  background: transparent;
  border: none;
}

.tooltip-popover strong {
  color: var(--vp-c-text-1);
  font-weight: 600;
}

@keyframes tooltipFadeIn {
  from {
    opacity: 0;
    transform: translateY(-4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.feature-info-compact {
  flex: 1;
  min-width: 0;
}

.feature-chip {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}
</style>
