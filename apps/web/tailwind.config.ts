import type { Config } from 'tailwindcss'
import dsPreset from '../../packages/ui-kit/tailwind.preset'

export default {
  presets: [dsPreset],
  content: [
    './components/**/*.{vue,ts}',
    './pages/**/*.{vue,ts}',
    './layouts/**/*.{vue,ts}',
    './composables/**/*.{vue,ts}',
    './plugins/**/*.{vue,ts}',
    '../../packages/ui-kit/src/**/*.{vue,ts}',
  ],
} satisfies Config
