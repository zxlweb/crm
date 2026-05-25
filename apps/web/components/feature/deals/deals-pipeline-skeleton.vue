<template>
  <div class="space-y-6" data-testid="deals-pipeline-skeleton">
    <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <div
        v-for="i in 4"
        :key="i"
        class="h-[108px] animate-pulse rounded-ds-xl border border-ds-border/60 bg-ds-bg-muted/50 motion-reduce:animate-none"
      />
    </div>
    <div class="flex gap-4 overflow-x-auto pb-2">
      <div
        v-for="col in 5"
        :key="col"
        class="flex w-72 shrink-0 flex-col rounded-xl border border-ds-border/60 bg-ds-bg-muted/30"
      >
        <div class="border-b border-ds-border/60 px-3 py-3">
          <div class="h-4 w-24 animate-pulse rounded bg-ds-bg-muted motion-reduce:animate-none" />
          <div class="mt-2 h-3 w-16 animate-pulse rounded bg-ds-bg-muted motion-reduce:animate-none" />
        </div>
        <div class="space-y-2 p-2">
          <div
            v-for="card in 2"
            :key="card"
            class="h-[120px] animate-pulse rounded-xl border border-ds-border/40 bg-ds-bg motion-reduce:animate-none"
          />
        </div>
      </div>
    </div>
  </div>
</template>
