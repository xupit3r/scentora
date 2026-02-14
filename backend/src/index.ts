import { createApp } from './app.js';
import { config } from './config.js';
import { prisma } from './prisma.js';
import { ensureSystemUser } from './services/systemUser.service.js';

const app = createApp();

await ensureSystemUser();

const server = app.listen(config.port, () => {
  console.log(`Scentora API running on port ${config.port}`);
});

// Graceful shutdown
async function shutdown() {
  console.log('Shutting down...');
  server.close();
  await prisma.$disconnect();
  process.exit(0);
}

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);
