import { defineConfig } from 'astro/config';
import preact from '@astrojs/preact';

// https://astro.build/config
export default defineConfig({
    base: '/astroterm/',
    integrations: [preact()],
    site: 'https://pkg.spooky.click'
});
