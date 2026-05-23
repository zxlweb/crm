<template>
  <section
    class="overflow-hidden rounded-2xl border border-ds-border/60 bg-ds-bg-elevated shadow-ds-sm"
    data-testid="dashboard-pipeline-tabs"
  >
    <div class="border-b border-ds-border-muted px-5 py-4 sm:px-6">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h3 class="text-sm font-semibold text-ds-fg-heading">{{ $t('dashboardPipelineTitle') }}</h3>
          <p class="mt-0.5 text-xs text-ds-fg-muted">{{ $t('dashboardPipelineHint') }}</p>
        </div>
        <div class="flex items-center gap-3">
          <UiTabs v-model="activeTab" :items="tabs" class="max-w-xs" />
          <NuxtLink
            :to="viewAllHref"
            class="shrink-0 text-xs font-semibold text-ds-fg-brand transition-colors hover:text-ds-brand-hover"
          >
            {{ $t('dashboardViewAll') }}
          </NuxtLink>
        </div>
      </div>
    </div>

    <div class="p-2 sm:p-3">
      <div v-if="activeTab === 'leads'">
        <div v-if="leads.length === 0" class="px-4 py-10 text-center text-sm text-ds-fg-muted">
          {{ $t('leadsEmpty') }}
        </div>
        <ul v-else class="divide-y divide-ds-border-muted">
          <li v-for="lead in leads" :key="lead.id">
            <NuxtLink
              :to="`/leads/${lead.id}`"
              class="group flex items-center gap-3 rounded-xl px-3 py-3 transition-colors hover:bg-ds-bg-muted/50"
            >
              <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-ds-brand-subtle text-xs font-bold text-ds-fg-brand">
                {{ initials(lead.title) }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate font-medium text-ds-fg-heading group-hover:text-ds-fg-brand">
                  {{ lead.title }}
                </p>
                <p class="mt-0.5 text-xs text-ds-fg-muted">
                  {{ $t('leadsColEngagement') }} {{ lead.engagement_score }}
                </p>
              </div>
              <div class="flex shrink-0 flex-col items-end gap-1.5 sm:flex-row sm:items-center">
                <CrmRelationshipHealthBadge variant="plain" muted-dot :health="lead.relationship_health" />
                <CrmLeadStatusBadge variant="plain" muted-dot :status="lead.status" />
              </div>
            </NuxtLink>
          </li>
        </ul>
      </div>

      <div v-else>
        <div v-if="accounts.length === 0" class="px-4 py-10 text-center text-sm text-ds-fg-muted">
          {{ $t('accountsEmpty') }}
        </div>
        <ul v-else class="divide-y divide-ds-border-muted">
          <li v-for="account in accounts" :key="account.id">
            <NuxtLink
              :to="`/accounts/${account.id}`"
              class="group flex items-center gap-3 rounded-xl px-3 py-3 transition-colors hover:bg-ds-bg-muted/50"
            >
              <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-ds-bg-muted text-xs font-bold text-ds-fg-muted">
                {{ initials(account.name) }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate font-medium text-ds-fg-heading group-hover:text-ds-fg-brand">
                  {{ account.name }}
                </p>
                <p class="mt-0.5 text-xs text-ds-fg-muted">
                  {{ account.industry || '—' }} · {{ $t('leadsColEngagement') }} {{ account.engagement_score }}
                </p>
              </div>
              <CrmRelationshipHealthBadge variant="plain" muted-dot :health="account.relationship_health" />
            </NuxtLink>
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { Account } from '~/types/account'
import type { Lead } from '~/types/lead'

defineProps<{
  leads: Lead[]
  accounts: Account[]
}>()

const { t } = useI18n()
const activeTab = ref<'leads' | 'accounts'>('leads')

const tabs = computed(() => [
  { id: 'leads', label: t('crmNavLeads') },
  { id: 'accounts', label: t('crmNavAccounts') },
])

const viewAllHref = computed(() => (activeTab.value === 'leads' ? '/leads' : '/accounts'))

function initials(name: string) {
  return name.trim().slice(0, 1).toUpperCase() || '?'
}
</script>
