<template>
  <section
    class="ds-pipeline-activity relative overflow-hidden rounded-2xl border border-ds-border-muted bg-ds-bg-elevated shadow-ds-sm"
    data-testid="dashboard-pipeline-tabs"
  >
    <span
      class="pointer-events-none absolute inset-x-0 top-0 h-0.5 bg-ds-info opacity-80"
      aria-hidden="true"
    />

    <header
      class="flex flex-col gap-3 border-b border-ds-border-muted px-4 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-5"
    >
      <div class="flex min-w-0 items-start gap-2.5">
        <span
          class="mt-0.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-lg bg-ds-info-subtle text-ds-info ring-1 ring-inset ring-ds-info/20"
          aria-hidden="true"
        >
          <UIcon name="i-heroicons-arrow-path" class="h-3.5 w-3.5" />
        </span>
        <div class="min-w-0">
          <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardPipelineTitle') }}</h3>
          <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dashboardPipelineHint') }}</p>
        </div>
      </div>
      <div class="flex shrink-0 items-center gap-2">
        <div
          class="ds-pipeline-activity__tabs inline-flex rounded-xl border border-ds-border-muted bg-ds-bg-muted/40 p-0.5"
          role="tablist"
        >
          <button
            v-for="tab in tabs"
            :key="tab.id"
            type="button"
            role="tab"
            :aria-selected="activeTab === tab.id"
            class="inline-flex cursor-pointer items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs font-medium transition-colors duration-200"
            :class="activeTab === tab.id
              ? 'bg-ds-bg-elevated text-ds-fg-heading shadow-ds-sm ring-1 ring-inset ring-ds-border-muted'
              : 'text-ds-fg-muted hover:text-ds-fg-heading'"
            @click="activeTab = tab.id"
          >
            <UIcon :name="tab.icon" class="h-3.5 w-3.5" aria-hidden="true" />
            <span>{{ tab.label }}</span>
            <span
              v-if="tab.count > 0"
              class="rounded-full bg-ds-bg-muted px-1.5 py-0 text-[10px] font-semibold tabular-nums text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted"
              :class="activeTab === tab.id ? '!bg-ds-info-subtle !text-ds-info !ring-ds-info/20' : ''"
            >
              {{ tab.count }}
            </span>
          </button>
        </div>
        <NuxtLink
          :to="viewAllHref"
          class="inline-flex cursor-pointer items-center gap-0.5 text-xs font-semibold text-ds-fg-brand transition-colors duration-200 hover:text-ds-brand-hover"
        >
          {{ $t('dashboardViewAll') }}
          <UIcon name="i-heroicons-arrow-right-20-solid" class="h-3 w-3" aria-hidden="true" />
        </NuxtLink>
      </div>
    </header>

    <div class="p-3 sm:p-4">
      <div
        v-if="activeItems.length === 0"
        class="rounded-xl border border-dashed border-ds-border bg-ds-bg-muted/20 px-6 py-10 text-center"
      >
        <div
          class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-ds-info-subtle text-ds-info"
          aria-hidden="true"
        >
          <UIcon name="i-heroicons-inbox" class="h-5 w-5" />
        </div>
        <p class="text-sm font-medium text-ds-fg-heading">
          {{ activeTab === 'leads' ? $t('leadsEmpty') : $t('accountsEmpty') }}
        </p>
      </div>

      <ul
        v-else
        class="grid gap-1.5 sm:gap-2 lg:grid-cols-2"
      >
        <li v-for="row in activeRows" :key="row.id">
          <NuxtLink
            :to="row.href"
            class="ds-pipeline-activity__row group flex cursor-pointer items-center gap-3 rounded-xl border border-ds-border-muted bg-ds-bg-elevated/70 px-3 py-2.5 transition-[border-color,background-color,box-shadow] duration-200 hover:border-ds-info/30 hover:bg-ds-bg-muted/40 hover:shadow-ds-sm"
            :data-testid="`dashboard-pipeline-row-${row.entityType}-${row.id}`"
          >
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-xs font-bold"
              :class="row.entityType === 'lead'
                ? 'bg-ds-brand-subtle text-ds-fg-brand'
                : 'bg-ds-bg-muted text-ds-fg-muted ring-1 ring-inset ring-ds-border-muted'"
            >
              {{ row.initials }}
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-1.5">
                <p class="truncate text-sm font-medium text-ds-fg-heading group-hover:text-ds-fg-brand">
                  {{ row.title }}
                </p>
              </div>
              <div class="mt-0.5 flex flex-wrap items-center gap-x-2 gap-y-0.5 text-[11px] text-ds-fg-muted">
                <span class="inline-flex items-center gap-0.5">
                  <UIcon name="i-heroicons-bolt" class="h-3 w-3" aria-hidden="true" />
                  {{ row.engagement }}
                </span>
                <span v-if="row.subtitle" class="truncate">{{ row.subtitle }}</span>
                <span class="inline-flex items-center gap-0.5">
                  <UIcon name="i-heroicons-clock" class="h-3 w-3" aria-hidden="true" />
                  {{ row.timeAgo }}
                </span>
              </div>
            </div>

            <div class="flex shrink-0 items-center gap-1.5">
              <CrmRelationshipHealthBadge
                variant="plain"
                muted-dot
                :health="row.health"
              />
              <CrmLeadStatusBadge
                v-if="row.entityType === 'lead' && row.status"
                variant="plain"
                muted-dot
                :status="row.status"
              />
              <UIcon
                name="i-heroicons-chevron-right-20-solid"
                class="h-3.5 w-3.5 text-ds-fg-subtle transition-transform duration-200 group-hover:translate-x-0.5 group-hover:text-ds-fg-brand"
                aria-hidden="true"
              />
            </div>
          </NuxtLink>
        </li>
      </ul>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Account } from '~/types/account'
import type { Lead, LeadStatus, RelationshipHealth } from '~/types/lead'

const props = defineProps<{
  leads: Lead[]
  accounts: Account[]
}>()

const { t, locale } = useI18n()
const activeTab = ref<'leads' | 'accounts'>('leads')

const tabs = computed(() => [
  {
    id: 'leads' as const,
    label: t('crmNavLeads'),
    icon: 'i-heroicons-light-bulb',
    count: props.leads.length,
  },
  {
    id: 'accounts' as const,
    label: t('crmNavAccounts'),
    icon: 'i-heroicons-building-office-2',
    count: props.accounts.length,
  },
])

const viewAllHref = computed(() => (activeTab.value === 'leads' ? '/leads' : '/accounts'))

const activeItems = computed(() =>
  activeTab.value === 'leads' ? props.leads : props.accounts,
)

interface ActivityRow {
  id: string
  entityType: 'lead' | 'account'
  title: string
  subtitle: string
  href: string
  initials: string
  engagement: string
  timeAgo: string
  health: RelationshipHealth
  status?: LeadStatus
}

const relativeFormatter = computed(() =>
  new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }),
)

function timeAgoLabel(iso: string): string {
  if (!iso) return ''
  const then = new Date(iso).getTime()
  if (Number.isNaN(then)) return ''
  const diffSec = Math.round((then - Date.now()) / 1000)
  const abs = Math.abs(diffSec)
  if (abs < 60) return relativeFormatter.value.format(Math.round(diffSec), 'second')
  if (abs < 3600) return relativeFormatter.value.format(Math.round(diffSec / 60), 'minute')
  if (abs < 86_400) return relativeFormatter.value.format(Math.round(diffSec / 3600), 'hour')
  if (abs < 2_592_000) return relativeFormatter.value.format(Math.round(diffSec / 86_400), 'day')
  if (abs < 31_536_000) return relativeFormatter.value.format(Math.round(diffSec / 2_592_000), 'month')
  return relativeFormatter.value.format(Math.round(diffSec / 31_536_000), 'year')
}

function initials(name: string): string {
  const trimmed = name.trim()
  if (!trimmed) return '?'
  return trimmed.slice(0, 1).toUpperCase()
}

const activeRows = computed<ActivityRow[]>(() => {
  if (activeTab.value === 'leads') {
    return props.leads.slice(0, 8).map((lead) => ({
      id: lead.id,
      entityType: 'lead' as const,
      title: lead.title,
      subtitle: '',
      href: `/leads/${lead.id}`,
      initials: initials(lead.title),
      engagement: t('dashboardActivityEngagement', { value: lead.engagement_score }),
      timeAgo: timeAgoLabel(lead.updated_at),
      health: lead.relationship_health,
      status: lead.status,
    }))
  }
  return props.accounts.slice(0, 8).map((account) => ({
    id: account.id,
    entityType: 'account' as const,
    title: account.name,
    subtitle: account.industry || '',
    href: `/accounts/${account.id}`,
    initials: initials(account.name),
    engagement: t('dashboardActivityEngagement', { value: account.engagement_score }),
    timeAgo: timeAgoLabel(account.updated_at),
    health: account.relationship_health,
  }))
})
</script>
