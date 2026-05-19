<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Play, Square, RotateCcw, Globe, ExternalLink, RefreshCw } from 'lucide-vue-next'
import type { DesktopStatus } from '@/types'

defineProps<{
  status: DesktopStatus
  loading: boolean
}>()

const emit = defineEmits<{
  start: []
  stop: []
  restart: []
  openWebUI: []
  openBrowser: []
  refresh: []
}>()
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <Button size="sm" :disabled="loading || status.running" @click="emit('start')">
      <Play class="w-4 h-4 mr-1.5" />
      启动
    </Button>
    <Button size="sm" variant="secondary" :disabled="loading || !status.running || status.attached" @click="emit('stop')">
      <Square class="w-4 h-4 mr-1.5" />
      停止
    </Button>
    <Button size="sm" variant="secondary" :disabled="loading || status.attached" @click="emit('restart')">
      <RotateCcw class="w-4 h-4 mr-1.5" />
      重启
    </Button>
    <Button size="sm" variant="outline" :disabled="loading" @click="emit('openWebUI')">
      <Globe class="w-4 h-4 mr-1.5" />
      打开 Web UI
    </Button>
    <Button size="sm" variant="outline" :disabled="loading" @click="emit('openBrowser')">
      <ExternalLink class="w-4 h-4 mr-1.5" />
      浏览器打开
    </Button>
    <Button size="sm" variant="ghost" :disabled="loading" @click="emit('refresh')">
      <RefreshCw class="w-4 h-4 mr-1.5" />
      刷新
    </Button>
  </div>
</template>
