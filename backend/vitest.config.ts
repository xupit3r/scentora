import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    pool: 'forks',
    fileParallelism: false,
    setupFiles: ['./tests/setup.ts'],
    include: ['tests/**/*.test.ts'],
    testTimeout: 30000,
    env: {
      JWT_SECRET: 'test-secret-key-for-vitest-do-not-use-in-production',
      DATABASE_URL: 'postgres://admin:password@localhost:5435/scentora?sslmode=disable',
      PORT: '3001',
      NODE_ENV: 'test',
    },
  },
});
