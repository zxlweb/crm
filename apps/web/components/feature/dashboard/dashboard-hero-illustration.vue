<template>
  <div
    class="dashboard-hero-illustration pointer-events-none select-none"
    aria-hidden="true"
  >
    <DashboardWorkbenchIllustration v-if="!isDark" class="dashboard-hero-illustration__svg" />
    <DashboardWorkbenchIllustrationDark v-else class="dashboard-hero-illustration__svg" />
  </div>
</template>

<script setup lang="ts">
const { isDark } = useTheme()
</script>

<style scoped>
.dashboard-hero-illustration :deep(.dashboard-workbench-illustration__svg),
.dashboard-hero-illustration__svg {
  width: 100%;
  max-width: 15rem;
  height: auto;
  filter: drop-shadow(0 20px 40px rgba(124, 58, 237, 0.12));
}

@media (min-width: 1280px) {
  .dashboard-hero-illustration :deep(.dashboard-workbench-illustration__svg),
  .dashboard-hero-illustration__svg {
    max-width: 17.5rem;
  }
}

@media (prefers-reduced-motion: no-preference) {
  .dashboard-hero-illustration {
    animation: dashboard-hero-float 7s ease-in-out infinite;
  }
}

@keyframes dashboard-hero-float {
  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-0.375rem);
  }
}

@media (prefers-reduced-motion: reduce) {
  .dashboard-hero-illustration {
    animation: none;
  }
}
</style>
