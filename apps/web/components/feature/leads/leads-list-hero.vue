<template>
  <section
    class="ds-leads-hero relative overflow-hidden rounded-2xl border border-ds-border bg-ds-bg-elevated/85 p-4 shadow-ds-sm backdrop-blur-sm sm:p-5"
    data-testid="leads-list-hero"
  >
    <span
      class="pointer-events-none absolute -right-12 -top-12 h-44 w-44 rounded-full opacity-25 blur-3xl"
      :style="{ background: 'var(--ds-brand-gradient)' }"
      aria-hidden="true"
    />

    <div class="relative grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      <article
        v-for="(stat, idx) in stats"
        :key="stat.key"
        class="ds-leads-hero__tile group relative overflow-hidden rounded-xl border border-ds-border-muted bg-ds-bg-elevated/70 p-3.5 transition-[border-color,box-shadow,transform] duration-200 hover:-translate-y-0.5 hover:border-ds-brand-muted/60 hover:shadow-ds-md"
        :class="stat.featured ? 'ds-leads-hero__tile--featured' : ''"
      >
        <span
          class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
          :class="stat.accentClass"
          aria-hidden="true"
        />
        <span
          v-if="stat.featured"
          class="pointer-events-none absolute -right-8 -top-8 h-24 w-24 rounded-full opacity-25 blur-2xl"
          :style="{ background: 'var(--ds-brand-gradient)' }"
          aria-hidden="true"
        />

        <div class="relative flex items-start gap-3">
          <span
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl transition-colors duration-200"
            :class="stat.iconWrapClass"
            :style="stat.iconStyle"
          >
            <UIcon :name="stat.icon" class="h-4 w-4" aria-hidden="true" />
          </span>
          <div class="min-w-0 flex-1">
            <div class="flex items-center justify-between gap-2">
              <p class="truncate text-[11px] font-medium uppercase tracking-wider text-ds-fg-muted">
                {{ stat.label }}
              </p>
              <span
                v-if="stat.tagLabel"
                class="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 text-[10px] font-semibold leading-none ring-1 ring-inset"
                :class="stat.tagClass"
              >
                <UIcon
                  v-if="stat.tagIcon"
                  :name="stat.tagIcon"
                  class="h-2.5 w-2.5"
                  aria-hidden="true"
                />
                {{ stat.tagLabel }}
              </span>
            </div>
            <p class="mt-1 truncate text-[1.5rem] font-extrabold tabular-nums tracking-tight text-ds-fg-heading sm:text-2xl">
              <component
                :is="stat.href ? 'NuxtLink' : 'span'"
                v-if="stat.featured"
                :to="stat.href"
                class="bg-clip-text"
                :style="brandGradientText"
              >
                {{ stat.value }}
              </component>
              <component
                :is="stat.href ? 'NuxtLink' : 'span'"
                v-else
                :to="stat.href"
                class="text-ds-fg-heading transition-colors duration-200"
                :class="stat.href ? 'cursor-pointer hover:text-ds-fg-brand' : ''"
              >
                {{ stat.value }}
              </component>
            </p>
            <p class="mt-0.5 text-xs text-ds-fg-subtle">{{ stat.hint }}</p>
          </div>
        </div>

        <NuxtLink
          v-if="stat.href"
          :to="stat.href"
          :aria-label="stat.label"
          class="absolute inset-0"
          :tabindex="idx === 0 ? 0 : -1"
        />
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Lead } from '~/types/lead'

interface HeroStat {
  key: string
  label: string
  value: string | number
  hint: string
  icon: string
  iconWrapClass: string
  iconStyle?: Record<string, string>
  accentClass: string
  featured?: boolean
  tagLabel?: string
  tagIcon?: string
  tagClass?: string
  href?: string
}

const props = defineProps<{
  items: Lead[]
  total: number
}>()

const { t } = useI18n()

const brandGradientText = {
  background: 'var(--ds-brand-gradient)',
  '-webkit-background-clip': 'text',
  'background-clip': 'text',
  color: 'transparent',
}

const HOT_THRESHOLD = 70
const IDLE_THRESHOLD_DAYS = 14

const hotCount = computed(
  () => props.items.filter((l) => l.engagement_score >= HOT_THRESHOLD).length,
)

const atRiskCount = computed(
  () => props.items.filter((l) => l.relationship_health === 'low').length,
)

const idleCount = computed(() => {
  const now = Date.now()
  return props.items.filter((l) => {
    if (!l.last_activity_at) return true
    const diff = Math.floor((now - new Date(l.last_activity_at).getTime()) / 86_400_000)
    return diff > IDLE_THRESHOLD_DAYS
  }).length
})

const stats = computed<HeroStat[]>(() => [
  {
    key: 'total',
    label: t('leadsHeroStatTotal'),
    value: props.total,
    hint: t('leadsHeroStatTotalHint', { n: props.items.length }),
    icon: 'i-heroicons-rectangle-stack',
    iconWrapClass: 'text-ds-on-brand shadow-ds-brand',
    iconStyle: { background: 'var(--ds-brand-gradient)' },
    accentClass: 'bg-ds-brand opacity-70',
    featured: true,
    tagLabel: t('leadsHeroStatLive'),
    tagIcon: 'i-heroicons-bolt',
    tagClass: 'bg-ds-brand-subtle text-ds-fg-brand ring-ds-brand/25',
  },
  {
    key: 'hot',
    label: t('leadsHeroStatHot'),
    value: hotCount.value,
    hint: t('leadsHeroStatHotHint'),
    icon: 'i-heroicons-fire',
    iconWrapClass: 'bg-ds-success-subtle text-ds-success',
    accentClass: 'bg-ds-success opacity-70',
    tagLabel:
      hotCount.value > 0
        ? t('leadsHeroStatHotTag')
        : t('leadsHeroStatQuietTag'),
    tagIcon: hotCount.value > 0 ? 'i-heroicons-arrow-trending-up' : 'i-heroicons-minus',
    tagClass:
      hotCount.value > 0
        ? 'bg-ds-success-subtle text-ds-success ring-ds-success/25'
        : 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted',
  },
  {
    key: 'at-risk',
    label: t('leadsHeroStatAtRisk'),
    value: atRiskCount.value,
    hint: t('leadsHeroStatAtRiskHint'),
    icon: 'i-heroicons-shield-exclamation',
    iconWrapClass:
      atRiskCount.value > 0
        ? 'bg-ds-danger-subtle text-ds-danger'
        : 'bg-ds-bg-muted text-ds-fg-muted',
    accentClass:
      atRiskCount.value > 0 ? 'bg-ds-danger opacity-70' : 'bg-ds-border opacity-50',
    tagLabel:
      atRiskCount.value > 0 ? t('leadsHeroStatActNow') : t('leadsHeroStatHealthy'),
    tagIcon:
      atRiskCount.value > 0
        ? 'i-heroicons-exclamation-triangle'
        : 'i-heroicons-check-circle',
    tagClass:
      atRiskCount.value > 0
        ? 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/25'
        : 'bg-ds-success-subtle text-ds-success ring-ds-success/25',
    href: atRiskCount.value > 0 ? '/leads?health=low' : undefined,
  },
  {
    key: 'idle',
    label: t('leadsHeroStatIdle'),
    value: idleCount.value,
    hint: t('leadsHeroStatIdleHint', { days: IDLE_THRESHOLD_DAYS }),
    icon: 'i-heroicons-clock',
    iconWrapClass:
      idleCount.value > 0
        ? 'bg-ds-warning-subtle text-ds-warning'
        : 'bg-ds-bg-muted text-ds-fg-muted',
    accentClass:
      idleCount.value > 0 ? 'bg-ds-warning opacity-70' : 'bg-ds-border opacity-50',
    tagLabel:
      idleCount.value > 0
        ? t('leadsHeroStatFollowUp')
        : t('leadsHeroStatOnTrack'),
    tagIcon:
      idleCount.value > 0
        ? 'i-heroicons-bell-alert'
        : 'i-heroicons-check-circle',
    tagClass:
      idleCount.value > 0
        ? 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/25'
        : 'bg-ds-success-subtle text-ds-success ring-ds-success/25',
  },
])
</script>

<style scoped>
.ds-leads-hero__tile--featured {
  background-image: linear-gradient(
    180deg,
    color-mix(in srgb, var(--ds-bg-elevated) 92%, var(--ds-brand) 8%) 0%,
    var(--ds-bg-elevated) 70%
  );
}
</style>
