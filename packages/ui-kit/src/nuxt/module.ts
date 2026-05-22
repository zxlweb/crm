import {
  addComponentsDir,
  addImports,
  addPlugin,
  createResolver,
  defineNuxtModule,
} from '@nuxt/kit'

export default defineNuxtModule({
  meta: {
    name: '@crm/ui-kit',
    configKey: 'uiKit',
    compatibility: { nuxt: '>=3.16.0' },
  },
  defaults: {},
  setup(_options, nuxt) {
    const resolver = createResolver(import.meta.url)
    const srcDir = resolver.resolve('..')

    nuxt.options.css.push(resolver.resolve('../styles/design-system.css'))

    addComponentsDir({
      path: resolver.resolve('../components/ui/chart'),
      pathPrefix: false,
      prefix: 'Chart',
      global: false,
    })
    addComponentsDir({
      path: resolver.resolve('../components/ui/card'),
      pathPrefix: false,
      prefix: 'Card',
      global: false,
    })
    addComponentsDir({
      path: resolver.resolve('../components/ui'),
      pathPrefix: false,
      prefix: 'Ui',
      ignore: ['chart', 'card'],
      global: false,
    })

    addPlugin(resolver.resolve('../runtime/echarts.client'))

    addImports([
      {
        name: 'useChartTheme',
        from: resolver.resolve('../chart/use-chart-theme'),
      },
    ])
  },
})
