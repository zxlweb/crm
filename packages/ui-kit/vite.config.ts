import { resolve } from 'node:path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    lib: {
      entry: {
        index: resolve(__dirname, 'src/index.ts'),
        'tokens/index': resolve(__dirname, 'src/tokens/index.ts'),
      },
      formats: ['es'],
    },
    rollupOptions: {
      external: [
        'vue',
        'echarts',
        'echarts/core',
        'echarts/charts',
        'echarts/components',
        'echarts/renderers',
        'vue-echarts',
      ],
      output: {
        assetFileNames: (info) => (info.name === 'style.css' ? 'ui-kit.css' : '[name][extname]'),
        entryFileNames: '[name].js',
      },
    },
    cssCodeSplit: false,
  },
})
