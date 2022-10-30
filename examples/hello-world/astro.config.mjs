// @ts-check
import { defineConfig } from 'astro/config';

export default defineConfig({
  vite: {
    server: {
      watch: {
        ignored: ['debug.log']
      }
    },
  }
});