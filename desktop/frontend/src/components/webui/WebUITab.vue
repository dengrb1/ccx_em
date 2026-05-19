<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Globe } from 'lucide-vue-next'
import type { DesktopStatus } from '@/types'
import { OpenWebUIInBrowser } from '@bindings/github.com/BenedictKing/ccx/desktop/desktopservice'

const props = defineProps<{
  status: DesktopStatus
  loading: boolean
}>()

const iframeSrc = computed(() => props.status.url)

const openInBrowser = async () => {
  try {
    await OpenWebUIInBrowser()
  } catch {
    // handled by parent
  }
}
</script>

<template>
  <div>
    <div v-if="status.running && iframeSrc" class="rounded-lg overflow-hidden border border-border" style="min-height: 620px">
      <iframe
        :src="iframeSrc"
        class="w-full border-0"
        style="min-height: 620px; background: white"
        title="CCX Web UI"
      />
    </div>
    <Card v-else>
      <CardContent class="flex flex-col items-start gap-4 py-8">
        <p class="text-sm text-muted-foreground">CCX 服务尚未启动，无法显示 Web UI。</p>
        <Button size="sm" :disabled="loading" @click="openInBrowser">
          <Globe class="w-4 h-4 mr-1.5" />
          浏览器打开
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
