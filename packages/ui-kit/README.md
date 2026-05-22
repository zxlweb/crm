# @crm/ui-kit

EnterpriseFlow 设计系统与图表组件库，可从 CRM monorepo 独立发布。

## 在 Nuxt 应用中使用

```ts
// nuxt.config.ts
export default defineNuxtConfig({
  modules: ['@crm/ui-kit/nuxt'],
})
```

```ts
// apps/web/plugins/00-ui-kit-theme.ts
import { bridgeUiKitThemeFromApp } from '@crm/ui-kit'

export default defineNuxtPlugin(() => {
  bridgeUiKitThemeFromApp(useTheme())
})
```

应用需实现 `useTheme()`（cookie + `applyThemeToDocument`）。ui-kit 通过模块级 bridge 消费主题，避免 monorepo 分包下 `inject` 失效。修改 `src/` 后请执行 `pnpm --filter @crm/ui-kit build`（或依赖 `exports.development` 指向 `src`）。

```vue
<UiThemeToggle :format-label="(id) => $t(id === 'v1' ? 'themeV1' : 'themeV2')" />
<ChartLine :categories="days" :series="series" loading-text="加载中…" />
```

## 开发

```bash
pnpm --filter @crm/ui-kit test
pnpm --filter @crm/ui-kit build
```

## 发布

```bash
pnpm --filter @crm/ui-kit publish --access restricted
```

构建产物：`dist/index.js`、`dist/ui-kit.css`、`dist/*.d.ts`。
