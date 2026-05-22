<template>
  <div class="space-y-10">
    <section
      v-for="layer in tokenLayers"
      :key="layer.id"
      class="ds-card rounded-2xl p-6"
    >
      <h3 class="text-ds-lg font-ds-semibold text-ds-fg-heading">
        {{ $t(layer.titleKey) }}
      </h3>
      <p class="mt-1 text-ds-sm text-ds-fg-muted">
        {{ $t(layer.descKey) }}
      </p>

      <div
        v-for="group in groupsForLayer(layer.id)"
        :key="group.id"
        class="mt-8"
      >
        <h4 class="text-ds-sm font-ds-medium uppercase tracking-wider text-ds-fg-muted">
          {{ $t(group.titleKey) }}
        </h4>
        <div class="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="entry in group.tokens"
            :key="entry.name + entry.token"
            class="rounded-ds-lg border border-ds-border bg-ds-bg-muted p-3"
          >
            <div class="mb-3 flex min-h-[3rem] items-center justify-center">
              <template v-if="entry.preview === 'color'">
                <div
                  class="h-10 w-full rounded-ds-md border border-ds-border"
                  :style="{ background: `var(${entry.token})` }"
                />
              </template>
              <template v-else-if="entry.preview === 'shadow'">
                <div
                  class="h-10 w-20 rounded-ds-md bg-ds-bg-elevated"
                  :style="{ boxShadow: `var(${entry.token})` }"
                />
              </template>
              <template v-else-if="entry.preview === 'radius'">
                <div
                  class="h-10 w-16 border-2 border-ds-brand bg-ds-bg-elevated"
                  :style="{ borderRadius: `var(${entry.token})` }"
                />
              </template>
              <template v-else-if="entry.preview === 'space'">
                <div
                  class="rounded-ds-sm bg-ds-brand"
                  :style="{ width: `var(${entry.token})`, height: `var(${entry.token})` }"
                />
              </template>
              <template v-else-if="entry.preview === 'text-sample'">
                <span
                  class="text-ds-fg-heading"
                  :style="textSampleStyle(entry.token)"
                >Aa</span>
              </template>
              <template v-else-if="entry.preview === 'font'">
                <span
                  class="text-ds-base text-ds-fg"
                  :style="{ fontFamily: `var(${entry.token})` }"
                >Ag</span>
              </template>
              <template v-else-if="entry.noteKey">
                <span class="text-ds-xs text-ds-fg-muted">{{ $t(entry.noteKey) }}</span>
              </template>
              <template v-else>
                <code class="text-ds-xs text-ds-fg-muted">{{ entry.token }}</code>
              </template>
            </div>
            <p class="text-ds-sm font-ds-medium text-ds-fg-heading">{{ entry.name }}</p>
            <p v-if="entry.token !== '—'" class="font-mono text-[10px] text-ds-fg-muted">{{ entry.token }}</p>
            <p v-if="entry.tailwind" class="mt-1 text-[10px] text-ds-fg-brand">
              {{ $t('dsTokenTailwind') }}: <code>{{ entry.tailwind }}</code>
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import {
  DS_TOKEN_GROUPS,
  DS_TOKEN_LAYERS,
  type DsTokenLayerId,
} from '@crm/ui-kit/tokens'

/** 供模板 v-for 使用 */
const tokenLayers = DS_TOKEN_LAYERS

function groupsForLayer(layer: DsTokenLayerId) {
  return DS_TOKEN_GROUPS.filter((g) => g.layer === layer)
}

function textSampleStyle(token: string): Record<string, string> {
  if (token.includes('font-normal') || token.includes('font-medium') || token.includes('font-semibold') || token.includes('font-bold')) {
    return { fontWeight: `var(${token})`, fontSize: 'var(--ds-text-xl)' }
  }
  const leading = token.replace('--ds-text-', '--ds-leading-')
  return { fontSize: `var(${token})`, lineHeight: `var(${leading})` }
}
</script>
