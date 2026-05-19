<script setup lang="ts">
import MetricsGrid from '@/components/status/MetricsGrid.vue'
import ServiceActions from '@/components/status/ServiceActions.vue'
import ServiceDetails from '@/components/status/ServiceDetails.vue'
import LogViewer from '@/components/status/LogViewer.vue'
import { useStatus } from '@/composables/useStatus'

const { status, loading, actionError, startService, stopService, restartService, openInBrowser, refresh } = useStatus()

const emit = defineEmits<{
  switchToWeb: []
}>()
</script>

<template>
  <div class="space-y-4">
    <MetricsGrid :status="status" />
    <ServiceActions
      :status="status"
      :loading="loading"
      @start="startService"
      @stop="stopService"
      @restart="restartService"
      @open-web-u-i="emit('switchToWeb')"
      @open-browser="openInBrowser"
      @refresh="refresh"
    />
    <p v-if="actionError" class="text-sm text-destructive-foreground">{{ actionError }}</p>
    <p v-else-if="status.lastError" class="text-sm text-destructive-foreground">{{ status.lastError }}</p>
    <ServiceDetails :status="status" />
    <LogViewer :logs="status.logs" />
  </div>
</template>
