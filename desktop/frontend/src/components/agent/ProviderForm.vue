<script setup lang="ts">
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import type { AgentProvider } from '@/types'

const props = defineProps<{
  selectedProvider: AgentProvider
  providerKeys: Record<AgentProvider, string>
  miMOBaseUrl: string
}>()

const emit = defineEmits<{
  'update:selectedProvider': [value: AgentProvider]
  'update:providerKeys': [value: Record<AgentProvider, string>]
  'update:miMOBaseUrl': [value: string]
}>()

const onProviderChange = (e: Event) => {
  emit('update:selectedProvider', (e.target as HTMLSelectElement).value as AgentProvider)
}

const onKeyChange = (value: string | number) => {
  emit('update:providerKeys', {
    ...props.providerKeys,
    [props.selectedProvider]: String(value),
  })
}
</script>

<template>
  <div class="space-y-3">
    <div class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">Provider</Label>
      <select
        :value="selectedProvider"
        class="w-full h-9 rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
        @change="onProviderChange"
      >
        <option value="ccx">CCX 本地网关</option>
        <option value="deepseek">DeepSeek 直连</option>
        <option value="mimo">MiMo 直连</option>
      </select>
    </div>

    <div v-if="selectedProvider !== 'ccx'" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">API Key</Label>
      <Input
        type="password"
        autocomplete="off"
        placeholder="仅写入 Claude Code 配置"
        :model-value="providerKeys[selectedProvider]"
        @update:model-value="onKeyChange"
      />
    </div>

    <div v-if="selectedProvider === 'mimo'" class="space-y-1.5">
      <Label class="text-xs text-muted-foreground">MiMo Base URL</Label>
      <Input
        type="url"
        placeholder="https://api.mimo.xiaomi.com/v1"
        :model-value="miMOBaseUrl"
        @update:model-value="emit('update:miMOBaseUrl', String($event))"
      />
    </div>
  </div>
</template>
