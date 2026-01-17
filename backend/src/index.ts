import Koa from 'koa';
import cors from '@koa/cors';
import bodyParser from 'koa-bodyparser';
import { config } from './config';
import { initDatabase } from './config/database';
import { errorHandler } from './middleware/errorHandler';
import router from './routes';

const app = new Koa();

// Middleware
app.use(errorHandler);
app.use(cors());
app.use(bodyParser());

// Routes
app.use(router.routes());
app.use(router.allowedMethods());

// Error event handler
app.on('error', (err, ctx) => {
  console.error('Server error:', err);
});

async function start() {
  try {
    // Initialize database
    await initDatabase();
    
    // Start server
    app.listen(config.port, () => {
      console.log(`🚀 Server running on http://localhost:${config.port}`);
      console.log(`📊 Environment: ${config.nodeEnv}`);
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

start();
