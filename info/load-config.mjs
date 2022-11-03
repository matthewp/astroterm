const path = require('path');
const url = require('url');

const main = require.resolve('astro');
const entryDir = path.dirname(main);

(async function() {
  const { openConfig } = await import(path.join(entryDir, '/dist/core/config/index.js'));

  const c = await openConfig(process.cwd(), undefined, 'dev', true, {}, false);
  const { astroConfig: config } = c;
  
  const parts = {
    srcDir: url.fileURLToPath(config.srcDir),
    outDir: url.fileURLToPath(config.outDir),
    base: config.base,
    output: config.output,
    build: {
      format: config.build.format,
      client: url.fileURLToPath(config.build.client),
      server: url.fileURLToPath(config.build.server),
      serverEntry: config.build.serverEntry,
    },
    integrations: config.integrations.map(int => int.name)
  };
  
  console.log(JSON.stringify(parts));
})();