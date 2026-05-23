<template>
  <div
    class="dashboard-page relative mx-auto max-w-[1400px] animate-pulse space-y-6 motion-reduce:animate-none"
    data-testid="dashboard-skeleton"
    aria-busy="true"
    aria-label="Loading dashboard"
  >
    <div class="pointer-events-none absolute inset-x-0 top-0 h-64 overflow-hidden" aria-hidden="true">
      <div class="absolute -left-16 top-0 h-48 w-48 rounded-full bg-ds-bg-muted/80 blur-3xl" />
      <div class="absolute right-0 top-8 h-36 w-36 rounded-full bg-ds-bg-muted/60 blur-3xl" />
    </div>

    <section>
      <div class="mb-4 space-y-2">
        <div class="h-3 w-24 rounded bg-ds-bg-muted" />
        <div class="h-7 w-48 max-w-full rounded bg-ds-bg-muted" />
      </div>
      <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div v-for="i in 4" :key="i" class="h-[6.5rem] rounded-ds-xl border border-ds-border bg-ds-bg-elevated shadow-ds-sm" />
      </div>
    </section>

    <div class="grid gap-4 xl:grid-cols-[minmax(0,1fr)_min(340px,100%)]">
      <section class="rounded-2xl border border-ds-border/60 bg-ds-bg-elevated">
        <div class="border-b border-ds-border-muted px-4 py-3 sm:px-5">
          <div class="h-4 w-24 rounded bg-ds-bg-muted" />
          <div class="mt-2 h-3 w-40 rounded bg-ds-bg-muted" />
        </div>
        <div class="grid gap-2 p-3 sm:p-4 xl:grid-cols-2">
          <div v-for="i in 4" :key="i" class="h-[5rem] rounded-xl bg-ds-bg-muted/70" />
        </div>
      </section>
      <aside class="min-h-[22rem] rounded-2xl border border-ds-border/60 bg-ds-bg-elevated xl:min-h-0" />
    </div>

    <div class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_min(360px,100%)]">
      <div class="h-64 rounded-2xl border border-ds-border/60 bg-ds-bg-elevated" />
      <aside class="h-72 rounded-2xl border border-ds-border/60 bg-ds-bg-elevated xl:h-auto xl:min-h-[18rem]" />
    </div>
  </div>
</template>
