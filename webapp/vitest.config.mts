/// <reference types="vitest" />
/// <reference types="vite/client" />

import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'
import tsconfigPaths from 'vite-tsconfig-paths'

// https://vitejs.dev/config/
export default defineConfig({
  test: {
    environment: 'happy-dom',
    setupFiles: 'test/vitest-setup-file.ts',
    include: [
      // Matches vitest tests in any subfolder of 'test'
      // having extension ~ /[.](test|spec)\.[cm]?[jt]s/
      'test/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts}'
    ]
  },
  plugins: [
    vue({
      template: { transformAssetUrls }
    }),
    quasar({
      sassVariables: 'src/quasar-variables.scss'
    }),
    tsconfigPaths()
  ]
})
