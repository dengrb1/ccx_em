<script setup lang="ts">
import MetricsGrid from '@/components/status/MetricsGrid.vue'
import ServiceActions from '@/components/status/ServiceActions.vue'
import ServiceDetails from '@/components/status/ServiceDetails.vue'
import DiagnosticCard from '@/components/status/DiagnosticCard.vue'
import LogViewer from '@/components/status/LogViewer.vue'
import { useStatus } from '@/composables/useStatus'
import { useResponsesDiagnostics } from '@/composables/useResponsesDiagnostics'
import { useLanguage } from '@/composables/useLanguage'

const { status, loading, actionError, startService, stopService, restartService, openInBrowser, refresh } = useStatus()
const { statusSummaryVisible, statusSummaryText } = useResponsesDiagnostics()
const { t } = useLanguage()

const emit = defineEmits<{
  switchToDashboard: []
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
      @open-web-u-i="emit('switchToDashboard')"
      @open-browser="openInBrowser"
      @refresh="refresh"
    />
    <DiagnosticCard
      v-if="actionError"
      :error="actionError"
      @dismiss="actionError = ''"
    />
    <DiagnosticCard
      v-else-if="status.lastError"
      :error="status.lastError"
    />
    <div
      v-if="!actionError && !status.lastError && statusSummaryVisible"
      class="rounded-lg border border-amber-500/20 bg-amber-500/5 backdrop-blur-sm px-4 py-3"
    >
      <div class="space-y-1">
        <h4 class="text-sm font-semibold text-amber-700 dark:text-amber-400">{{ t('agent.codexDiagnosticStatusTitle') }}</h4>
        <p class="text-xs text-muted-foreground leading-relaxed">{{ statusSummaryText }}</p>
        <p class="text-xs text-muted-foreground">{{ t('agent.codexDiagnosticStatusHint') }}</p>
      </div>
    </div>
    <ServiceDetails :status="status" />
    <LogViewer :logs="status.logs" />
  </div>
</template>
