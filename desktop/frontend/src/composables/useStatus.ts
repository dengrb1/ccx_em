import { ref, onMounted, onBeforeUnmount } from 'vue'
import type { DesktopStatus } from '@/types'
import {
  GetStatus,
  StartService,
  StopService,
  RestartService,
  OpenWebUIInBrowser,
} from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'

// Module-level singletons — all composables share the same state
const status = ref<DesktopStatus>({
  running: false,
  starting: false,
  attached: false,
  port: 0,
  url: '',
  pid: 0,
  binaryPath: '',
  dataDir: '',
  logs: [],
})
const loading = ref(false)
const actionError = ref('')
let statusInterval: ReturnType<typeof setInterval> | undefined

const syncStatus = async () => {
  try {
    const data = (await GetStatus()) as DesktopStatus
    status.value = {
      ...status.value,
      ...data,
      logs: Array.isArray(data.logs) ? data.logs : [],
    }
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : String(error)
  }
}

const invoke = async (action: () => Promise<unknown>) => {
  actionError.value = ''
  try {
    await action()
    await syncStatus()
  } catch (error) {
    actionError.value = error instanceof Error ? error.message : String(error)
  }
}

const startService = () => invoke(StartService)
const stopService = () => invoke(StopService)
const restartService = () => invoke(RestartService)
const openInBrowser = () => invoke(OpenWebUIInBrowser)

const refresh = async () => {
  loading.value = true
  try {
    await syncStatus()
  } finally {
    loading.value = false
  }
}

export function useStatus() {
  onMounted(async () => {
    await syncStatus()
    if (!statusInterval) {
      statusInterval = setInterval(syncStatus, 3000)
    }
  })

  onBeforeUnmount(() => {
    if (statusInterval) {
      clearInterval(statusInterval)
      statusInterval = undefined
    }
  })

  return {
    status,
    loading,
    actionError,
    syncStatus,
    startService,
    stopService,
    restartService,
    openInBrowser,
    refresh,
  }
}
