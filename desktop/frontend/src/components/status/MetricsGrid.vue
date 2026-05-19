<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import type { DesktopStatus } from '@/types'

defineProps<{
  status: DesktopStatus
}>()
</script>

<template>
  <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
    <Card v-for="item in [
      { label: '端口', value: status.port },
      { label: '版本', value: status.health?.version?.version || 'v0.0.0-dev' },
      { label: '运行时长', value: status.health?.uptime ? `${Math.floor(status.health.uptime / 60)}m` : '--' },
      { label: '上游数', value: status.health?.config?.upstreamCount || 0 },
    ]" :key="item.label">
      <CardContent class="p-4">
        <p class="text-xs text-muted-foreground mb-2">{{ item.label }}</p>
        <p class="text-lg font-semibold">{{ item.value }}</p>
      </CardContent>
    </Card>
  </div>
</template>
