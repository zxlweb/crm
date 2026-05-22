<template>
  <div class="min-h-screen bg-ds-bg">
    <header class="sticky top-0 z-10 border-b border-ds-border bg-ds-bg-topbar px-6 py-4 backdrop-blur-md">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4">
        <div>
          <p class="text-xs font-medium text-ds-fg-brand">{{ $t('cardsPageLabel') }}</p>
          <h1 class="text-xl font-bold text-ds-fg-heading">{{ $t('cardsPageTitle') }}</h1>
          <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('cardsPageDesc') }}</p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <UiThemeToggle />
          <NuxtLink to="/charts" class="text-sm text-ds-fg-muted transition-colors hover:text-ds-fg-brand">
            {{ $t('chartsPageTitle') }}
          </NuxtLink>
          <NuxtLink to="/design" class="text-sm text-ds-fg-muted transition-colors hover:text-ds-fg-brand">
            {{ $t('themePreviewTitle') }}
          </NuxtLink>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl space-y-10 px-6 py-8">
      <section>
        <h2 class="mb-1 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">
          {{ $t('cardScenarioDashboard') }}
        </h2>
        <p class="mb-4 text-sm text-ds-fg-muted">{{ $t('cardsSectionDashboardDesc') }}</p>
        <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <CardMetric
            :label="$t('cardsDemoOutput')"
            :value="22"
            :compare-label="$t('metricCompareYoy')"
            trend="12%"
            trend-direction="down"
            icon-tone="info"
          >
            <template #icon>
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </template>
          </CardMetric>

          <CardMetric
            :label="$t('tenantActive')"
            :value="1284"
            :compare-label="$t('metricCompareActiveRate')"
            trend="+8.2%"
            trend-direction="up"
            icon-tone="calendar"
          >
            <template #icon>
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </template>
          </CardMetric>

          <CardMetric
            :label="$t('cardsDemoRevenue')"
            :value="88692"
            :compare-label="$t('metricCompareYoy')"
            trend="+4.1%"
            trend-direction="up"
            icon-tone="accent"
          >
            <template #icon>
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </template>
          </CardMetric>

          <CardMetric
            :label="$t('tenantTotal')"
            :value="12"
            :compare-label="$t('metricCompareTotal')"
            trend="—"
            trend-direction="flat"
            icon-tone="brand"
          >
            <template #icon>
              <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5" />
              </svg>
            </template>
          </CardMetric>
        </div>
      </section>

      <section>
        <h2 class="mb-1 text-sm font-semibold uppercase tracking-wider text-ds-fg-muted">
          {{ $t('cardScenarioContent') }}
        </h2>
        <p class="mb-4 text-sm text-ds-fg-muted">{{ $t('cardsSectionContentDesc') }}</p>
        <CardShell :title="$t('cardsShellTitle')" :subtitle="$t('cardsShellSubtitle')">
          <p class="text-sm text-ds-fg-muted">{{ $t('cardsShellBody') }}</p>
        </CardShell>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth' })
</script>
