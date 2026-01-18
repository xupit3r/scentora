import dotenv from 'dotenv';

dotenv.config();

export const config = {
  port: process.env.PORT || 3000,
  nodeEnv: process.env.NODE_ENV || 'development',
  jwtSecret: process.env.JWT_SECRET || 'your-secret-key-change-in-production',
  jwtAccessExpiresIn: process.env.JWT_ACCESS_EXPIRES_IN || '15m',
  jwtRefreshExpiresIn: process.env.JWT_REFRESH_EXPIRES_IN || '7d',
  couchdb: {
    url: process.env.COUCHDB_URL || 'http://localhost:5984',
    user: process.env.COUCHDB_USER || 'admin',
    password: process.env.COUCHDB_PASSWORD || 'password',
    database: process.env.COUCHDB_DATABASE || 'scentora',
  },
};
