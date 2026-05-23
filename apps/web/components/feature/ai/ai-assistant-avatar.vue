<template>
  <div
    class="ai-assistant-avatar"
    :class="sizeClass"
    role="img"
    :aria-label="$t('aiAssistantAvatarAria')"
  >
    <span v-if="!minimal" class="ai-assistant-avatar__pulse" aria-hidden="true" />
    <span v-if="!minimal" class="ai-assistant-avatar__orbit" aria-hidden="true">
      <span class="ai-assistant-avatar__spark ai-assistant-avatar__spark--a" />
      <span class="ai-assistant-avatar__spark ai-assistant-avatar__spark--b" />
    </span>

    <svg
      class="ai-assistant-avatar__svg"
      viewBox="0 0 64 64"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <defs>
        <linearGradient :id="gradBody" x1="16" y1="12" x2="48" y2="52" gradientUnits="userSpaceOnUse">
          <stop stop-color="#A78BFA" />
          <stop offset="1" stop-color="#6D28D9" />
        </linearGradient>
        <linearGradient :id="gradFace" x1="22" y1="24" x2="42" y2="44" gradientUnits="userSpaceOnUse">
          <stop stop-color="#F5F3FF" />
          <stop offset="1" stop-color="#EDE9FE" />
        </linearGradient>
        <radialGradient :id="gradGlow" cx="0" cy="0" r="1" gradientUnits="userSpaceOnUse" gradientTransform="translate(32 30) rotate(90) scale(24)">
          <stop stop-color="#C4B5FD" stop-opacity="0.55" />
          <stop offset="1" stop-color="#C4B5FD" stop-opacity="0" />
        </radialGradient>
      </defs>

      <circle cx="32" cy="34" r="24" :fill="`url(#${gradGlow})`" class="ai-assistant-avatar__inner-glow" />

      <!-- 天线 -->
      <line x1="32" y1="10" x2="32" y2="18" stroke="#8B5CF6" stroke-width="2" stroke-linecap="round" />
      <circle cx="32" cy="8" r="3" fill="#DDD6FE" class="ai-assistant-avatar__antenna" />
      <circle cx="32" cy="8" r="1.5" fill="#7C3AED" class="ai-assistant-avatar__antenna-core" />

      <!-- 主体 -->
      <circle cx="32" cy="36" r="22" :fill="`url(#${gradBody})`" />
      <circle cx="32" cy="36" r="22" stroke="#DDD6FE" stroke-width="1" stroke-opacity="0.45" />

      <!-- 面罩 -->
      <rect x="18" y="26" width="28" height="18" rx="9" :fill="`url(#${gradFace})`" />
      <rect x="18" y="26" width="28" height="18" rx="9" stroke="#C4B5FD" stroke-width="0.75" stroke-opacity="0.6" />

      <!-- 眼睛 -->
      <g class="ai-assistant-avatar__eye ai-assistant-avatar__eye--left">
        <ellipse cx="26" cy="35" rx="3" ry="3.5" fill="#5B21B6" />
        <circle cx="27" cy="34" r="1" fill="white" opacity="0.9" />
      </g>
      <g class="ai-assistant-avatar__eye ai-assistant-avatar__eye--right">
        <ellipse cx="38" cy="35" rx="3" ry="3.5" fill="#5B21B6" />
        <circle cx="39" cy="34" r="1" fill="white" opacity="0.9" />
      </g>

      <!-- 微笑 -->
      <path
        d="M27 40 Q32 43.5 37 40"
        stroke="#7C3AED"
        stroke-width="1.5"
        stroke-linecap="round"
        fill="none"
      />

      <!-- 耳侧光点 -->
      <circle cx="12" cy="36" r="2" fill="#C4B5FD" opacity="0.65" class="ai-assistant-avatar__cheek ai-assistant-avatar__cheek--left" />
      <circle cx="52" cy="36" r="2" fill="#C4B5FD" opacity="0.65" class="ai-assistant-avatar__cheek ai-assistant-avatar__cheek--right" />
    </svg>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    size?: 'xs' | 'sm' | 'md' | 'lg'
    /** 内联场景（徽章等）隐藏光晕与轨道粒子 */
    minimal?: boolean
  }>(),
  { size: 'md', minimal: false },
)

const sizeClass = computed(() => [
  `ai-assistant-avatar--${props.size}`,
  props.minimal ? 'ai-assistant-avatar--minimal' : '',
])

const uid = useId().replace(/:/g, '')
const gradBody = `ai-avatar-body-${uid}`
const gradFace = `ai-avatar-face-${uid}`
const gradGlow = `ai-avatar-glow-${uid}`
</script>

<style scoped>
.ai-assistant-avatar {
  position: relative;
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
}

.ai-assistant-avatar--xs {
  width: 1.125rem;
  height: 1.125rem;
}

.ai-assistant-avatar--sm {
  width: 2.25rem;
  height: 2.25rem;
}

.ai-assistant-avatar--md {
  width: 3rem;
  height: 3rem;
}

.ai-assistant-avatar--lg {
  width: 4rem;
  height: 4rem;
}

.ai-assistant-avatar__svg {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 6px 14px rgba(124, 58, 237, 0.22));
}

.ai-assistant-avatar__pulse {
  position: absolute;
  inset: 0;
  border-radius: 9999px;
  background: radial-gradient(circle, rgba(167, 139, 250, 0.35) 0%, rgba(167, 139, 250, 0) 70%);
}

.ai-assistant-avatar__orbit {
  position: absolute;
  inset: -0.25rem;
  z-index: 0;
}

.ai-assistant-avatar__spark {
  position: absolute;
  width: 0.35rem;
  height: 0.35rem;
  border-radius: 9999px;
  background: #c4b5fd;
  box-shadow: 0 0 6px rgba(196, 181, 253, 0.9);
}

.ai-assistant-avatar__spark--a {
  top: 0.15rem;
  left: 50%;
  transform: translateX(-50%);
}

.ai-assistant-avatar__spark--b {
  bottom: 0.35rem;
  right: 0.1rem;
}

.ai-assistant-avatar--minimal .ai-assistant-avatar__svg {
  filter: none;
}

@media (prefers-reduced-motion: no-preference) {
  .ai-assistant-avatar__svg {
    animation: ai-assistant-float 4.5s ease-in-out infinite;
  }

  .ai-assistant-avatar__pulse {
    animation: ai-assistant-pulse 2.8s ease-in-out infinite;
  }

  .ai-assistant-avatar__orbit {
    animation: ai-assistant-orbit 9s linear infinite;
  }

  .ai-assistant-avatar__spark--b {
    animation: ai-assistant-spark 2.2s ease-in-out infinite;
  }

  .ai-assistant-avatar__eye {
    transform-box: fill-box;
    animation: ai-assistant-blink 4.8s ease-in-out infinite;
  }

  .ai-assistant-avatar__eye--left {
    transform-origin: 26px 35px;
  }

  .ai-assistant-avatar__eye--right {
    transform-origin: 38px 35px;
    animation-delay: 0.08s;
  }

  .ai-assistant-avatar__antenna-core {
    animation: ai-assistant-antenna 1.6s ease-in-out infinite;
  }

  .ai-assistant-avatar__inner-glow {
    animation: ai-assistant-glow 3.2s ease-in-out infinite;
  }

  .ai-assistant-avatar__cheek--left {
    animation: ai-assistant-cheek 3s ease-in-out infinite;
  }

  .ai-assistant-avatar__cheek--right {
    animation: ai-assistant-cheek 3s ease-in-out infinite 0.4s;
  }
}

@keyframes ai-assistant-float {
  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-0.2rem);
  }
}

@keyframes ai-assistant-pulse {
  0%,
  100% {
    opacity: 0.45;
    transform: scale(0.92);
  }

  50% {
    opacity: 0.85;
    transform: scale(1.08);
  }
}

@keyframes ai-assistant-orbit {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@keyframes ai-assistant-spark {
  0%,
  100% {
    opacity: 0.35;
    transform: scale(0.85);
  }

  50% {
    opacity: 1;
    transform: scale(1.15);
  }
}

@keyframes ai-assistant-blink {
  0%,
  44%,
  48%,
  100% {
    transform: scaleY(1);
  }

  46% {
    transform: scaleY(0.12);
  }
}

@keyframes ai-assistant-antenna {
  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.45;
  }
}

@keyframes ai-assistant-glow {
  0%,
  100% {
    opacity: 0.55;
  }

  50% {
    opacity: 0.9;
  }
}

@keyframes ai-assistant-cheek {
  0%,
  100% {
    opacity: 0.45;
  }

  50% {
    opacity: 0.85;
  }
}

@media (prefers-reduced-motion: reduce) {
  .ai-assistant-avatar__svg,
  .ai-assistant-avatar__pulse,
  .ai-assistant-avatar__orbit,
  .ai-assistant-avatar__spark--b,
  .ai-assistant-avatar__eye,
  .ai-assistant-avatar__antenna-core,
  .ai-assistant-avatar__inner-glow,
  .ai-assistant-avatar__cheek {
    animation: none !important;
  }
}
</style>
