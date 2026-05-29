<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: number | string
  animated?: boolean
}>(), {
  size: 32,
  animated: true
})

const sizeStyle = computed(() => {
  const s = typeof props.size === 'number' ? `${props.size}px` : props.size
  return { width: s, height: s }
})
</script>

<template>
  <div class="ccx-logo-container" :style="sizeStyle">
    <svg
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="ccx-logo-svg"
      aria-hidden="true"
    >
      <defs>
        <!-- Web 导航专用轻量流光符号 -->
        <linearGradient id="web-logo-flow" x1="18%" y1="24%" x2="84%" y2="76%">
          <stop offset="0%" stop-color="#38bdf8" />
          <stop offset="48%" stop-color="#6366f1" />
          <stop offset="100%" stop-color="#10b981" />
        </linearGradient>
        <filter id="web-logo-glow" x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur stdDeviation="2.4" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      <!-- 轻量版 Terminal Gateway：无 App 背景，适配 Web 顶栏小尺寸 -->
      <g filter="url(#web-logo-glow)" stroke="url(#web-logo-flow)" stroke-linecap="round" stroke-linejoin="round">
        <path d="M 22 29 L 43 50 L 22 71" stroke-width="8.5" />
        <path d="M 56 28 L 80 72" stroke-width="8.5" />
        <path d="M 80 28 L 56 72" stroke-width="8.5" />
      </g>

      <!-- 细状态线和在线节点，呼应 App 图标但不压重导航栏 -->
      <path d="M 24 82 H 50" stroke="#10b981" stroke-width="3.2" stroke-linecap="round" opacity="0.55" />
      <path d="M 56 82 H 70" stroke="#38bdf8" stroke-width="3.2" stroke-linecap="round" opacity="0.42" />
      <g :class="{ 'animate-gateway-pulse': animated }">
        <circle cx="80" cy="82" r="3" fill="#5eead4" />
        <circle cx="80" cy="82" r="7" stroke="#5eead4" stroke-width="1.4" opacity="0.24" />
      </g>
    </svg>
  </div>
</template>

<style scoped>
.ccx-logo-container {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ccx-logo-svg {
  width: 100%;
  height: 100%;
  overflow: visible;
}

/* Web 导航在线节点呼吸脉冲 */
@keyframes gateway-pulse {
  0%, 100% {
    transform: scale(0.9);
    transform-origin: 80px 82px;
    opacity: 0.78;
  }
  50% {
    transform: scale(1.14);
    transform-origin: 80px 82px;
    opacity: 1;
  }
}

.animate-gateway-pulse {
  animation: gateway-pulse 2.4s infinite ease-in-out;
}
</style>
