<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'

defineProps<{
  logs: string[]
}>()
</script>

<template>
  <Card>
    <CardHeader class="pb-3">
      <CardTitle class="text-sm font-medium text-muted-foreground">最近日志</CardTitle>
    </CardHeader>
    <CardContent>
      <ScrollArea class="h-[260px]">
        <pre class="text-xs text-muted-foreground whitespace-pre-wrap break-all leading-relaxed">{{ logs.length ? logs.join('\n') : '暂无日志' }}</pre>
      </ScrollArea>
    </CardContent>
  </Card>
</template>
