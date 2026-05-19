<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { DesktopStatus } from '@/types'

defineProps<{
  status: DesktopStatus
}>()
</script>

<template>
  <Card>
    <CardHeader class="pb-3">
      <CardTitle class="text-sm font-medium text-muted-foreground">服务详情</CardTitle>
    </CardHeader>
    <CardContent class="space-y-3">
      <div v-for="item in [
        { label: '二进制', value: status.binaryPath || '未发现' },
        { label: '数据目录', value: status.dataDir || '未设置' },
        { label: 'PID', value: String(status.pid || '-') },
        { label: '健康状态', value: status.health?.status || 'unknown' },
      ]" :key="item.label" class="flex items-center justify-between text-sm">
        <span class="text-muted-foreground">{{ item.label }}</span>
        <code class="text-xs bg-secondary px-2 py-1 rounded-md break-all max-w-[60%] text-right">{{ item.value }}</code>
      </div>
    </CardContent>
  </Card>
</template>
