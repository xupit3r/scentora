import { PrismaClient } from '@prisma/client';
import { CLIContext } from './types/context.js';
import dotenv from 'dotenv';

// Load environment variables
dotenv.config();

let context: CLIContext | null = null;

/**
 * Initialize and return application context
 */
export async function initContext(): Promise<CLIContext> {
  if (context) {
    return context;
  }

  const prisma = new PrismaClient();

  // Test database connection
  try {
    await prisma.$connect();
  } catch (error) {
    console.error('❌ Failed to connect to database');
    console.error('Error:', error instanceof Error ? error.message : error);
    console.error('\nPlease check:');
    console.error('1. PostgreSQL is running (docker compose up -d)');
    console.error('2. DATABASE_URL in .env is correct');
    process.exit(1);
  }

  context = {
    prisma,
    config: {
      databaseUrl: process.env.DATABASE_URL || '',
      environment: process.env.NODE_ENV || 'development',
    },
  };

  return context;
}

/**
 * Get current context (must be initialized first)
 */
export function getContext(): CLIContext {
  if (!context) {
    throw new Error('Context not initialized. Call initContext() first.');
  }
  return context;
}

/**
 * Cleanup context and close connections
 */
export async function closeContext(): Promise<void> {
  if (context) {
    await context.prisma.$disconnect();
    context = null;
  }
}
