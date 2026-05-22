<template>
  <div class="space-y-6" data-testid="leads-reports-panel">
    <p class="text-sm text-ds-fg-muted">{{ $t('leadsReportsHint') }}</p>
    <section class="grid gap-6 lg:grid-cols-2">
      <CardShell :title="$t('leadsReportStatus')" :subtitle="$t('leadsReportStatusHint')" :height="280">
        <div class="flex h-[240px] items-center justify-center rounded-xl border border-dashed border-ds-border text-sm text-ds-fg-muted">
          {{ $t('leadsReportPlaceholder') }}
          <span class="ml-1 font-mono text-xs">ChartBar</span>
        </div>
      </CardShell>
      <CardShell :title="$t('leadsReportSource')" :subtitle="$t('leadsReportSourceHint')" :height="280">
        <div class="flex h-[240px] items-center justify-center rounded-xl border border-dashed border-ds-border text-sm text-ds-fg-muted">
          {{ $t('leadsReportPlaceholder') }}
          <span class="ml-1 font-mono text-xs">ChartDonut</span>
        </div>
      </CardShell>
    </section>
    <section class="grid gap-6 lg:grid-cols-2">
      <CardShell :title="$t('leadsReportTrend')" :subtitle="$t('leadsReportTrendHint')" :height="300">
        <div class="flex h-[260px] items-center justify-center rounded-xl border border-dashed border-ds-border text-sm text-ds-fg-muted">
          {{ $t('leadsReportPlaceholder') }}
          <span class="ml-1 font-mono text-xs">ChartLine</span>
        </div>
      </CardShell>
      <CardShell :title="$t('leadsReportFunnel')" :subtitle="$t('leadsReportFunnelHint')" :height="300">
        <div class="flex h-[260px] items-center justify-center rounded-xl border border-dashed border-ds-border text-sm text-ds-fg-muted">
          {{ $t('leadsReportPlaceholder') }}
          <span class="ml-1 font-mono text-xs">ChartFunnel</span>
        </div>
      </CardShell>
    </section>
  </div>
</template>
