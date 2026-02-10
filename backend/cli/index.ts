#!/usr/bin/env node
import Vorpal from 'vorpal';
import chalk from 'chalk';
import { initContext, closeContext } from './context.js';
import { registerUserCommands } from './commands/user.js';

const VERSION = '1.0.0';

/**
 * Main CLI entry point
 */
async function main() {
  console.log(chalk.bold.magenta('🌸 Scentora CLI Console') + chalk.gray(` (v${VERSION})`));
  console.log(chalk.gray('━'.repeat(40)));

  // Initialize context
  const context = await initContext();
  console.log(chalk.green('✅ Connected to database'));
  console.log(chalk.gray(`   Environment: ${context.config.environment}`));
  console.log(chalk.gray(`   Database: ${extractDbInfo(context.config.databaseUrl)}`));
  console.log('');
  console.log(chalk.dim("Type 'help' for available commands"));
  console.log(chalk.dim("Type 'exit' to quit"));
  console.log('');

  // Create Vorpal instance
  const vorpal = new Vorpal();

  // Register command modules
  registerUserCommands(vorpal);

  // Set up delimiter
  vorpal.delimiter(chalk.cyan('scentora> ')).show();

  // Handle exit
  vorpal.on('client_command_executed', (command: string) => {
    if (command === 'exit') {
      closeContext()
        .then(() => {
          console.log(chalk.gray('\nGoodbye! 👋'));
          process.exit(0);
        })
        .catch(err => {
          console.error('Error closing context:', err);
          process.exit(1);
        });
    }
  });

  // Handle Ctrl+C
  process.on('SIGINT', async () => {
    console.log('');
    await closeContext();
    console.log(chalk.gray('Goodbye! 👋'));
    process.exit(0);
  });
}

/**
 * Extract readable database info from connection string
 */
function extractDbInfo(url: string): string {
  try {
    const parsed = new URL(url);
    return `${parsed.hostname}:${parsed.port || '5432'}${parsed.pathname}`;
  } catch {
    return 'Unknown';
  }
}

// Run main
main().catch(error => {
  console.error(chalk.red('\n❌ Fatal error:'), error);
  process.exit(1);
});
