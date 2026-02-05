import 'dotenv/config';

export const config = {
  port: parseInt(process.env.PORT || '3000', 10),
  nodeEnv: process.env.NODE_ENV || 'development',

  databaseUrl: process.env.DATABASE_URL || '',

  jwtSecret: process.env.JWT_SECRET || '',
  jwtAccessExpiresIn: process.env.JWT_ACCESS_EXPIRES_IN || '15m',
  jwtRefreshExpiresIn: process.env.JWT_REFRESH_EXPIRES_IN || '7d',

  corsAllowedOrigins: (process.env.CORS_ALLOWED_ORIGINS || 'http://localhost:5173').split(','),

  rateLimitRequests: parseInt(process.env.RATE_LIMIT_REQUESTS || '100', 10),
  rateLimitWindow: process.env.RATE_LIMIT_WINDOW || '1m',
};

if (!config.jwtSecret) {
  throw new Error('JWT_SECRET environment variable is required');
}
