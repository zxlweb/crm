<template>
  <!-- 头部元信息行：轻量 inline -->
  <div
    v-if="variant === 'inline'"
    class="inline-flex min-w-0 items-center gap-2 text-xs text-ds-fg-muted sm:text-sm"
    data-testid="lead-owner-chip"
  >
    <CrmUserAvatar :name="profile.name" :seed="profile.id" size="sm" tone="neutral" />
    <span>{{ $t('leadsFieldOwner') }}</span>
    <span class="font-medium text-ds-fg">{{ profile.name }}</span>
    <span v-if="profile.role" class="hidden sm:inline">· {{ profile.role }}</span>
  </div>

  <!-- 卡片式（列表等场景备用） -->
  <div
    v-else
    class="inline-flex items-center gap-3 rounded-xl border border-ds-border bg-ds-bg-muted/60 px-3 py-2"
    data-testid="lead-owner-chip"
  >
    <CrmUserAvatar :name="profile.name" :seed="profile.id" size="md" tone="neutral" />
    <div class="min-w-0">
      <p class="text-xs font-medium text-ds-fg-muted">{{ $t('leadsFieldOwner') }}</p>
      <p class="truncate text-sm font-semibold text-ds-fg-heading">{{ profile.name }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { OwnerProfile } from '~/composables/use-owner-profile'

withDefaults(
  defineProps<{
    profile: OwnerProfile
    variant?: 'inline' | 'card'
  }>(),
  { variant: 'inline' },
)
</script>
