<template>
  <PermissionGuard resource="deals" action="view">
    <div class="space-y-4" data-testid="deal-detail-page">
      <NuxtLink to="/deals" class="text-sm text-ds-fg-brand hover:underline">{{ $t('dealsBackToList') }}</NuxtLink>

      <div v-if="pending" class="flex justify-center py-24">
        <UIcon name="i-heroicons-arrow-path" class="h-8 w-8 animate-spin text-primary" />
      </div>

      <UAlert v-else-if="loadError || !deal" color="red" variant="soft" :title="loadError || $t('dealsNotFound')" />

      <CardShell v-else :title="deal.title" :subtitle="dealStageLabel(deal.stage)">
        <dl class="grid gap-4 sm:grid-cols-2">
<div>
            <dt class="text-xs text-ds-fg-muted">{{ $t('dealsFieldAmount') }}</dt>
            <dd class="text-sm font-medium tabular-nums">{{ formatDealAmount(deal.amount, deal.currency) }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ds-fg-muted">{{ $t('dealsFieldProbability') }}</dt>
            <dd class="text-sm">{{ deal.probability }}%</dd>
          </div>
          <p class="text-sm text-ds-fg-muted sm:col-span-2">{{ deal.description || '—' }}</p>
        </dl>
      </CardShell>
    </div>
  </PermissionGuard>
</template>

<script setup lang="ts">
import type { Deal } from '~/types/deal'

definePageMeta({ layout: 'app', middleware: 'auth' })

const route = useRoute()
const { t } = useI18n()
const dealsApi = useDeals()
const { dealStageLabel, formatDealAmount } = useDealLabels()

const deal = ref<Deal | null>(null)
const pending = ref(true)
const loadError = ref('')

async function load() {
  pending.value = true
  loadError.value = ''
  try {
    deal.value = await dealsApi.fetchById(String(route.params.id))
  } catch (e) {
    loadError.value = e instanceof Error ? e.message : t('loadFailed')
  } finally {
    pending.value = false
  }
}

onMounted(load)
watch(() => route.params.id, load)
</script>
