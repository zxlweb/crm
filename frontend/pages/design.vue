<template>
  <div class="min-h-screen bg-ds-bg px-6 py-12 text-ds-fg">
    <div class="mx-auto max-w-5xl">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p class="text-sm font-medium text-ds-fg-brand">{{ $t('themePreviewLabel') }}</p>
          <h1 class="mt-2 text-3xl font-bold tracking-tight text-ds-fg-heading">{{ $t('themePreviewTitle') }}</h1>
          <p class="mt-3 max-w-2xl text-ds-fg-muted">{{ $t('themePreviewDesc') }}</p>
        </div>
        <UiThemeToggle />
      </div>

      <p class="mt-4 text-sm text-ds-fg-muted">
        {{ $t('themeCurrent') }}: <strong class="text-ds-fg-heading">{{ id === 'v1' ? $t('themeV1Name') : $t('themeV2Name') }}</strong>
      </p>

      <div class="mt-10 grid gap-8 lg:grid-cols-2">
        <article
          class="overflow-hidden rounded-2xl border-2 transition-all"
          :class="id === 'v1' ? 'border-ds-brand shadow-ds-brand' : 'border-ds-border'"
        >
          <div class="border-b border-ds-border bg-ds-bg-muted px-6 py-4">
            <span class="rounded-full bg-ds-brand-subtle px-3 py-1 text-xs font-semibold text-ds-fg-brand">V1</span>
            <h2 class="mt-2 text-xl font-semibold text-ds-fg-heading">{{ $t('themeV1Name') }}</h2>
            <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('themeV1Desc') }}</p>
          </div>
          <div class="space-y-3 bg-ds-bg-elevated p-6">
            <div class="flex gap-2">
              <span class="h-3 w-3 rounded-full bg-violet-300" />
              <span class="h-3 w-3 rounded-full bg-purple-400" />
              <span class="h-3 w-3 rounded-full bg-fuchsia-300" />
              <span class="ml-2 text-xs text-ds-fg-muted">{{ $t('themeV1Palette') }}</span>
            </div>
            <ul class="space-y-2 text-sm text-ds-fg">
              <li>· {{ $t('themeV1Point1') }}</li>
              <li>· {{ $t('themeV1Point2') }}</li>
              <li>· {{ $t('themeV1Point3') }}</li>
            </ul>
            <div class="flex flex-wrap gap-3 pt-4">
              <button
                type="button"
                class="ds-btn-primary cursor-pointer rounded-xl px-4 py-2.5 text-sm font-medium"
                @click="applyAndGo('v1', '/login')"
              >
                {{ $t('themeApplyV1') }}
              </button>
              <NuxtLink
                to="/login"
                class="rounded-xl border border-ds-border px-4 py-2.5 text-sm font-medium text-ds-fg-brand transition-colors hover:bg-ds-bg-muted"
                @click="setTheme('v1')"
              >
                {{ $t('themePreviewLogin') }}
              </NuxtLink>
              <NuxtLink
                to="/admin"
                class="rounded-xl border border-ds-border px-4 py-2.5 text-sm font-medium text-ds-fg-muted transition-colors hover:bg-ds-bg-muted"
                @click="setTheme('v1')"
              >
                {{ $t('themePreviewAdmin') }}
              </NuxtLink>
            </div>
          </div>
        </article>

        <article
          class="overflow-hidden rounded-2xl border-2 transition-all"
          :class="id === 'v2' ? 'border-ds-brand shadow-ds-brand' : 'border-ds-border'"
        >
          <div class="border-b border-ds-border bg-ds-bg-muted px-6 py-4">
            <span class="rounded-full bg-ds-brand-subtle px-3 py-1 text-xs font-semibold text-ds-fg-brand">V2</span>
            <h2 class="mt-2 text-xl font-semibold text-ds-fg-heading">{{ $t('themeV2Name') }}</h2>
            <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('themeV2Desc') }}</p>
          </div>
          <div class="space-y-3 bg-ds-bg-elevated p-6">
            <div class="flex gap-2">
              <span class="h-3 w-3 rounded-full bg-[#0D0D0D] ring-1 ring-ds-border" />
              <span class="h-3 w-3 rounded-full bg-purple-500" />
              <span class="h-3 w-3 rounded-full bg-emerald-500" />
              <span class="ml-2 text-xs text-ds-fg-muted">{{ $t('themeV2Palette') }}</span>
            </div>
            <ul class="space-y-2 text-sm text-ds-fg">
              <li>· {{ $t('themeV2Point1') }}</li>
              <li>· {{ $t('themeV2Point2') }}</li>
              <li>· {{ $t('themeV2Point3') }}</li>
            </ul>
            <div class="flex flex-wrap gap-3 pt-4">
              <button
                type="button"
                class="ds-btn-primary cursor-pointer rounded-xl px-4 py-2.5 text-sm font-medium"
                @click="applyAndGo('v2', '/login')"
              >
                {{ $t('themeApplyV2') }}
              </button>
              <NuxtLink
                to="/login"
                class="rounded-xl border border-ds-border px-4 py-2.5 text-sm font-medium text-ds-fg-brand transition-colors hover:bg-ds-bg-muted"
                @click="setTheme('v2')"
              >
                {{ $t('themePreviewLogin') }}
              </NuxtLink>
              <NuxtLink
                to="/admin"
                class="rounded-xl border border-ds-border px-4 py-2.5 text-sm font-medium text-ds-fg-muted transition-colors hover:bg-ds-bg-muted"
                @click="setTheme('v2')"
              >
                {{ $t('themePreviewAdmin') }}
              </NuxtLink>
            </div>
          </div>
        </article>
      </div>

      <section class="ds-card mt-10 rounded-2xl p-6">
        <h3 class="font-semibold text-ds-fg-heading">{{ $t('themeTokensTitle') }}</h3>
        <p class="mt-1 text-sm text-ds-fg-muted">{{ $t('themeTokensDesc') }}</p>
        <div class="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
          <div v-for="swatch in swatches" :key="swatch.label" class="rounded-lg border border-ds-border p-3">
            <div class="h-10 rounded-md" :style="{ background: swatch.var }" />
            <p class="mt-2 text-xs font-medium text-ds-fg-heading">{{ swatch.label }}</p>
            <p class="text-[10px] text-ds-fg-muted">{{ swatch.token }}</p>
          </div>
        </div>
      </section>

      <p class="mt-10 text-center text-sm text-ds-fg-muted">{{ $t('themePreviewHint') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ThemeId } from '~/design-system/tokens'

definePageMeta({ layout: 'auth' })

const { id, setTheme } = useTheme()

const swatches = [
  { label: 'Background', token: '--ds-bg', var: 'var(--ds-bg)' },
  { label: 'Surface', token: '--ds-bg-elevated', var: 'var(--ds-bg-elevated)' },
  { label: 'Brand', token: '--ds-brand', var: 'var(--ds-brand)' },
  { label: 'Heading', token: '--ds-fg-heading', var: 'var(--ds-fg-heading)' },
]

function applyAndGo(theme: ThemeId, path: string) {
  setTheme(theme)
  navigateTo(path)
}
</script>
