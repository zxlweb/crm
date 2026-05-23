<template>
  <section class="overflow-hidden rounded-2xl border border-ds-border/60 bg-ds-bg-elevated shadow-ds-sm">
    <div class="flex flex-wrap items-start justify-between gap-3 border-b border-ds-border-muted px-5 py-4 sm:px-6">
      <div>
        <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardFunnelTitle') }}</h3>
        <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dashboardFunnelHint') }}</p>
      </div>
      <NuxtLink
        to="/leads?tab=reports"
        class="cursor-pointer text-xs font-semibold text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
      >
        {{ $t('dashboardPreviewSampleCta') }}
      </NuxtLink>
    </div>
    <div class="px-5 py-4 sm:px-6">
      <ChartFunnel :items="funnelItems" :height="220" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { DASHBOARD_FUNNEL_FIXTURE } from '~/fixtures/dashboard-preview'

const { t } = useI18n()

const funnelItems = computed(() =>
  DASHBOARD_FUNNEL_FIXTURE.map((item) => ({
    name: t(`dashboardFunnelStage.${item.name}`),
    value: item.value,
  })),
)
</script>
