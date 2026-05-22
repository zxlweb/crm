/**
 * 组件分层 ESLint（与 tests/component-layers.test.ts 互补）
 *
 * 安装：cd frontend && npm install
 * 运行：npm run lint
 */
import boundaries from 'eslint-plugin-boundaries'
import pluginVue from 'eslint-plugin-vue'
import vueParser from 'vue-eslint-parser'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['.nuxt/**', 'dist/**', 'node_modules/**'] },
  {
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
        extraFileExtensions: ['.vue'],
        sourceType: 'module',
      },
    },
    plugins: {
      boundaries,
      vue: pluginVue,
    },
    settings: {
      'boundaries/elements': [
        { type: 'base', pattern: 'components/base/**', mode: 'full' },
        { type: 'ui', pattern: 'components/ui/**', mode: 'full' },
        { type: 'common', pattern: 'components/common/**', mode: 'full' },
        { type: 'layout', pattern: 'components/layout/**', mode: 'full' },
        { type: 'feature', pattern: 'components/feature/**', mode: 'full' },
        { type: 'pages', pattern: 'pages/**', mode: 'full' },
        { type: 'composables', pattern: 'composables/**', mode: 'full' },
      ],
      'boundaries/include': ['components/**', 'pages/**', 'composables/**'],
    },
    rules: {
      'boundaries/element-types': [
        'error',
        {
          default: 'disallow',
          rules: [
            { from: ['base'], allow: ['base'] },
            { from: ['ui'], allow: ['base', 'ui'] },
            { from: ['common'], allow: ['base', 'ui', 'common'] },
            { from: ['layout'], allow: ['base', 'ui', 'common', 'layout'] },
            { from: ['feature'], allow: ['base', 'ui', 'common', 'layout'] },
            {
              from: ['pages'],
              allow: ['base', 'ui', 'common', 'layout', 'feature', 'composables'],
            },
            { from: ['composables'], allow: [] },
          ],
        },
      ],
    },
  },
)
