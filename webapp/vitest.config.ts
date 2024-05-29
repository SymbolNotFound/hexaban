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
      // Matches vitest tests in any subfolder of 'test'.
      'test/**/*.{test,spec}.{js,mjs,ts,mts}'
    ]
  },
  plugins: [
    vue({
      template: { transformAssetUrls }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    }) as any,
    quasar({
      sassVariables: 'src/quasar-variables.scss'
    }),
    tsconfigPaths()
  ]
})
