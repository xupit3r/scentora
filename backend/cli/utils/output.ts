import chalk from 'chalk';
import Table from 'cli-table3';
import ora, { Ora } from 'ora';

/**
 * Output utilities for formatted CLI messages
 */

export const output = {
  success: (message: string) => {
    console.log(chalk.green('✅ ' + message));
  },

  error: (message: string) => {
    console.log(chalk.red('❌ ' + message));
  },

  warning: (message: string) => {
    console.log(chalk.yellow('⚠️  ' + message));
  },

  info: (message: string) => {
    console.log(chalk.blue('ℹ️  ' + message));
  },

  section: (title: string) => {
    console.log('\n' + chalk.bold.cyan(title));
    console.log(chalk.gray('─'.repeat(40)));
  },

  newline: () => {
    console.log('');
  },
};

/**
 * Create a spinner for long-running operations
 */
export function spinner(text: string): Ora {
  return ora({
    text,
    color: 'cyan',
  }).start();
}

/**
 * Create a formatted table
 */
export function createTable(head: string[]): Table.Table {
  return new Table({
    head: head.map(h => chalk.bold.cyan(h)),
    style: {
      head: [],
      border: ['gray'],
    },
  });
}

/**
 * Format a date for display
 */
export function formatDate(date: Date | string): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return d.toISOString().replace('T', ' ').substring(0, 19);
}

/**
 * Truncate text with ellipsis
 */
export function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength - 3) + '...';
}
