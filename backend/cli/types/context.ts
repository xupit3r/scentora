import { PrismaClient } from '@prisma/client';

/**
 * CLI Application Context
 * Provides access to database, services, and configuration
 */
export interface CLIContext {
  prisma: PrismaClient;
  config: {
    databaseUrl: string;
    environment: string;
  };
}
