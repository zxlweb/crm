<template>
  <UiTable
    :rows="items"
    :columns="columns"
    :empty-state="{ label: $t('leadsEmpty') }"
    data-testid="leads-list-table"
  >
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <template #title-data="{ row }">
      <NuxtLink
        :to="`/leads/${row.id}`"
        class="ds-leads-row__title group/title flex min-w-0 cursor-pointer items-center gap-2.5"
      >
        <span
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-xs font-bold tracking-tight text-ds-fg-brand ring-1 ring-inset ring-ds-brand/15 transition-[background-color,box-shadow,transform] duration-200 group-hover/title:scale-[1.04] group-hover/title:shadow-ds-sm"
          :style="avatarStyle(row.id)"
        >
          {{ initialsOf(row.title) }}
        </span>
        <span class="flex min-w-0 flex-col">
          <span
            class="truncate text-sm font-semibold text-ds-fg-heading transition-colors duration-200 group-hover/title:text-ds-fg-brand"
          >
            {{ row.title }}
          </span>
          <span
            v-if="metaTags(row).length"
            class="mt-0.5 flex flex-wrap items-center gap-1"
          >
            <span
              v-for="tag in metaTags(row)"
              :key="tag.id"
              class="inline-flex items-center gap-0.5 rounded border px-1.5 py-0 text-[10px] font-medium leading-4"
              :class="tag.toneClass"
            >
              <UIcon
                v-if="tag.icon"
                :name="tag.icon"
                class="h-2.5 w-2.5"
                aria-hidden="true"
              />
              {{ tag.label }}
            </span>
          </span>
        </span>
      </NuxtLink>
    </template>

    <template #status-data="{ row }">
      <CrmLeadStatusBadge variant="filled" :status="row.status" />
    </template>

    <template #lifecycle_stage-data="{ row }">
      <CrmLifecycleBadge variant="filled" :stage="row.lifecycle_stage" />
    </template>

    <template #relationship_health-data="{ row }">
      <CrmRelationshipHealthBadge variant="filled" :health="row.relationship_health" />
    </template>

    <template #engagement_score-data="{ row }">
      <div class="ds-leads-row__engagement flex items-center justify-end gap-2">
        <div
          class="relative h-1.5 w-16 overflow-hidden rounded-full bg-ds-bg-muted"
          aria-hidden="true"
        >
          <span
            class="absolute inset-y-0 left-0 rounded-full transition-[width,background-color] duration-500"
            :class="engagementBarClass(row.engagement_score)"
            :style="{ width: `${Math.max(4, Math.min(100, row.engagement_score))}%` }"
          />
        </div>
        <span
          class="min-w-[2.25rem] text-right text-sm font-semibold tabular-nums"
          :class="engagementTextClass(row.engagement_score)"
        >
          {{ row.engagement_score }}
        </span>
      </div>
    </template>

    <template #source-data="{ row }">
      <span
        class="inline-flex items-center gap-1 rounded-md border border-ds-border bg-ds-bg-muted px-2 py-0.5 text-[11px] font-medium text-ds-fg-muted"
      >
        <UIcon name="i-heroicons-globe-alt" class="h-3 w-3" aria-hidden="true" />
        {{ formatSource(row.source) }}
      </span>
    </template>

    <template #owner-data="{ row }">
      <div class="flex items-center gap-2">
        <template v-if="row.owner_id">
          <CrmUserAvatar
            :name="ownerName(row.owner_id)"
            :seed="row.owner_id"
            size="sm"
            tone="accent"
          />
          <span class="min-w-0">
            <span class="block truncate text-sm font-medium text-ds-fg">
              {{ ownerName(row.owner_id) }}
            </span>
            <span
              v-if="isMine(row.owner_id)"
              class="text-[10px] font-semibold uppercase tracking-wider text-ds-fg-brand"
            >
              {{ $t('leadsPoolBadgeMine') }}
            </span>
            <span
              v-else
              class="text-[10px] uppercase tracking-wider text-ds-fg-subtle"
            >
              {{ $t('leadsPoolBadgeOthers') }}
            </span>
          </span>
        </template>
        <template v-else>
          <span
            class="inline-flex h-7 w-7 items-center justify-center rounded-full bg-ds-warning-subtle text-ds-warning ring-1 ring-inset ring-ds-warning/30"
            aria-hidden="true"
          >
            <UIcon name="i-heroicons-globe-asia-australia" class="h-3.5 w-3.5" />
          </span>
          <span class="flex flex-col">
            <span class="text-sm font-medium text-ds-warning">
              {{ $t('leadsPoolBadgePublic') }}
            </span>
            <span class="text-[10px] uppercase tracking-wider text-ds-fg-subtle">
              {{ $t('leadsPoolBadgePublicHint') }}
            </span>
          </span>
        </template>
      </div>
    </template>

    <template #recycle_in-data="{ row }">
      <span
        v-if="!row.owner_id"
        class="text-xs text-ds-fg-subtle"
        :aria-label="$t('leadsPoolRecycleNA')"
      >—</span>
      <span
        v-else-if="!poolSettings || !poolSettings.enabled"
        class="text-xs text-ds-fg-subtle"
      >—</span>
      <span
        v-else
        class="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-[11px] font-semibold ring-1 ring-inset"
        :class="recycleToneClass(daysToRecycle(row))"
      >
        <UIcon
          name="i-heroicons-clock"
          class="h-3 w-3"
          aria-hidden="true"
        />
        {{ recycleLabel(daysToRecycle(row)) }}
      </span>
    </template>

    <template #actions-data="{ row }">
      <div class="flex items-center justify-end gap-0.5">
        <button
          v-if="canClaim(row)"
          type="button"
          class="ds-leads-row-action inline-flex h-7 cursor-pointer items-center gap-1 rounded-lg bg-ds-brand-subtle px-2 text-[11px] font-semibold text-ds-fg-brand ring-1 ring-inset ring-ds-brand/20 transition-colors duration-150 hover:bg-ds-brand hover:text-ds-on-brand disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="pendingId === row.id"
          :data-testid="`leads-row-claim-${row.id}`"
          :aria-label="$t('leadsPoolActionClaim')"
          @click="$emit('claim', row)"
        >
          <UIcon name="i-heroicons-hand-raised" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ $t('leadsPoolActionClaim') }}</span>
        </button>
        <button
          v-if="canRelease(row)"
          type="button"
          class="ds-leads-row-action inline-flex h-7 cursor-pointer items-center gap-1 rounded-lg px-2 text-[11px] font-medium text-ds-fg-muted ring-1 ring-inset ring-ds-border transition-colors duration-150 hover:bg-ds-bg-muted hover:text-ds-fg disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="pendingId === row.id"
          :data-testid="`leads-row-release-${row.id}`"
          @click="$emit('release', row)"
        >
          <UIcon name="i-heroicons-arrow-uturn-left" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ $t('leadsPoolActionRelease') }}</span>
        </button>
        <button
          v-if="canTransfer(row)"
          type="button"
          class="ds-leads-row-action inline-flex h-7 cursor-pointer items-center gap-1 rounded-lg px-2 text-[11px] font-medium text-ds-fg-muted ring-1 ring-inset ring-ds-border transition-colors duration-150 hover:bg-ds-bg-muted hover:text-ds-fg disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="pendingId === row.id"
          :data-testid="`leads-row-transfer-${row.id}`"
          @click="$emit('transfer', row)"
        >
          <UIcon name="i-heroicons-arrow-right-circle" class="h-3.5 w-3.5" aria-hidden="true" />
          <span>{{ $t('leadsPoolActionTransfer') }}</span>
        </button>
        <CrmTableIconAction
          :to="`/leads/${row.id}`"
          icon="i-heroicons-eye-20-solid"
          :label="$t('leadsViewDetail')"
          data-testid="leads-row-view"
        />
      </div>
    </template>

    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </UiTable>
</template>

<script setup lang="ts">
import type { Lead, LeadPool, LeadPoolSettings } from '~/types/lead'

// 与 contacts/accounts list-table 复用同一 UiTableColumn 形状；
// @crm/ui-kit 在 dev 模式下通过 condition exports 解析，TS 缓存偶发
// 无法解析 `import type { UiTableColumn } from '@crm/ui-kit'`，
// 这里用本地类型别名规避，结构与 @crm/ui-kit/components/ui/table-types.ts 保持同构。
type UiTableColumn = { key: string; label: string; sortable?: boolean; class?: string }

interface MetaTag {
  id: string
  label: string
  icon?: string
  toneClass: string
}

const LOW_INTENT_TAG = '低意向'

const props = withDefaults(
  defineProps<{
    items: Lead[]
    /** 当前所在客户池视图，影响列展示与可用操作 */
    pool?: LeadPool
    /** 当前用户 id，用于判断「我的」 */
    currentUserId?: string | null
    /** 是否拥有领取权限（私海容量未达上限、有 claim 权限） */
    canClaimPool?: boolean
    /** 是否拥有释放/转交权限 */
    canManagePool?: boolean
    /** 客户池设置：用于显示自动回收倒计时 */
    poolSettings?: LeadPoolSettings | null
    /** 正在处理的行 id，用于禁用该行按钮 */
    pendingId?: string | null
  }>(),
  {
    pool: 'all',
    currentUserId: null,
    canClaimPool: true,
    canManagePool: true,
    poolSettings: null,
    pendingId: null,
  },
)

defineEmits<{
  claim: [row: Lead]
  release: [row: Lead]
  transfer: [row: Lead]
}>()

const { t } = useI18n()
const { leadSourceLabel } = useLeadLabels()
const ownerProfile = useOwnerProfile()

const columns = computed<UiTableColumn[]>(() => {
  const base: UiTableColumn[] = [
    { key: 'title', label: t('leadsColTitle'), sortable: true },
    { key: 'status', label: t('status'), sortable: true },
    { key: 'lifecycle_stage', label: t('leadsColLifecycle'), sortable: true },
    { key: 'relationship_health', label: t('leadsColHealth'), sortable: true },
    { key: 'engagement_score', label: t('leadsColEngagement'), sortable: true, class: 'text-right' },
    { key: 'source', label: t('leadsColSource'), sortable: true },
    { key: 'owner', label: t('leadsColOwner') },
  ]
  if (props.pool !== 'public') {
    base.push({ key: 'recycle_in', label: t('leadsColRecycleIn'), class: 'text-right whitespace-nowrap' })
  }
  base.push({ key: 'actions', label: t('actions'), class: 'text-right w-40' })
  return base
})

function ownerName(ownerId: string): string {
  return ownerProfile.resolve(ownerId)?.name ?? ownerId.slice(0, 8)
}

function isMine(ownerId: string | null): boolean {
  return Boolean(props.currentUserId && ownerId && ownerId === props.currentUserId)
}

function canClaim(row: Lead): boolean {
  return !row.owner_id && Boolean(props.canClaimPool)
}

function canRelease(row: Lead): boolean {
  if (!props.canManagePool) return false
  if (!row.owner_id) return false
  if (props.pool === 'mine') return isMine(row.owner_id)
  return isMine(row.owner_id)
}

function canTransfer(row: Lead): boolean {
  if (!props.canManagePool) return false
  if (!row.owner_id) return false
  return isMine(row.owner_id)
}

function daysToRecycle(row: Lead): number | null {
  if (!props.poolSettings || !props.poolSettings.enabled) return null
  if (!row.owner_id) return null
  const baselineRaw = row.last_activity_at ?? row.created_at
  const baseline = new Date(baselineRaw).getTime()
  const idleDays = (Date.now() - baseline) / 86_400_000
  return Math.ceil(props.poolSettings.inactiveDays - idleDays)
}

function recycleLabel(days: number | null): string {
  if (days == null) return '—'
  if (days <= 0) return t('leadsPoolRecycleOverdue')
  return t('leadsPoolRecycleInDays', { days })
}

function recycleToneClass(days: number | null): string {
  if (days == null) return 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
  if (days <= 0) return 'bg-ds-danger-subtle text-ds-danger ring-ds-danger/30'
  if (days <= 3) return 'bg-ds-warning-subtle text-ds-warning ring-ds-warning/30'
  if (days <= 7) return 'bg-ds-info-subtle text-ds-info ring-ds-info/20'
  return 'bg-ds-bg-muted text-ds-fg-muted ring-ds-border-muted'
}

function formatSource(source: string) {
  return leadSourceLabel(source) || '—'
}

function initialsOf(name: string): string {
  if (!name) return '·'
  const trimmed = name.trim()
  if (!trimmed) return '·'
  const ascii = /^[\u0020-\u007E]+$/.test(trimmed)
  if (ascii) {
    const words = trimmed.split(/\s+/).filter(Boolean).slice(0, 2)
    return words.map((w) => w[0]?.toUpperCase() ?? '').join('') || trimmed[0]!.toUpperCase()
  }
  return Array.from(trimmed)[0] ?? '·'
}

function hashCode(key: string): number {
  let h = 0
  for (let i = 0; i < key.length; i += 1) {
    h = Math.trunc((h << 5) - h + (key.codePointAt(i) ?? 0))
  }
  return Math.abs(h)
}

const AVATAR_TINTS = [
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-brand) 22%, transparent), color-mix(in srgb, var(--ds-brand) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-info) 22%, transparent), color-mix(in srgb, var(--ds-info) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-success) 22%, transparent), color-mix(in srgb, var(--ds-success) 8%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-warning) 24%, transparent), color-mix(in srgb, var(--ds-warning) 10%, transparent))',
  'linear-gradient(135deg, color-mix(in srgb, var(--ds-danger) 20%, transparent), color-mix(in srgb, var(--ds-danger) 8%, transparent))',
]

function avatarStyle(id: string) {
  const idx = hashCode(id) % AVATAR_TINTS.length
  return { background: AVATAR_TINTS[idx] }
}

function metaTags(row: Lead): MetaTag[] {
  const tags: MetaTag[] = []
  const hasLowIntent = (row.tags ?? []).some(
    (tag) => tag === LOW_INTENT_TAG || tag.toLowerCase() === 'low intent',
  )
  if (hasLowIntent) {
    tags.push({
      id: 'low-intent',
      label: t('leadsTagLowIntent'),
      icon: 'i-heroicons-eye-slash',
      toneClass: 'border-ds-border-muted bg-ds-bg-muted text-ds-fg-muted',
    })
  }
  if (row.amount > 0) {
    tags.push({
      id: 'amount-set',
      label: t('leadsRowTagAmountSet'),
      icon: 'i-heroicons-banknotes',
      toneClass: 'border-ds-info/25 bg-ds-info-subtle text-ds-info',
    })
  }
  return tags
}

function engagementBarClass(score: number): string {
  if (score >= 70) return 'bg-ds-success'
  if (score >= 40) return 'bg-ds-warning'
  return 'bg-ds-danger'
}

function engagementTextClass(score: number): string {
  if (score >= 70) return 'text-ds-success'
  if (score >= 40) return 'text-ds-fg'
  return 'text-ds-danger'
}
</script>
