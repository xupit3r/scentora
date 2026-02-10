import Vorpal from 'vorpal';
import inquirer from 'inquirer';
import bcrypt from 'bcryptjs';
import { getContext } from '../context.js';
import { output, spinner, createTable, formatDate } from '../utils/output.js';
import {
  isValidEmail,
  isValidPassword,
  validationMessages,
  isValidUUID,
} from '../utils/validation.js';

/**
 * Register user management commands
 */
export function registerUserCommands(vorpal: Vorpal) {
  // Create user command
  vorpal
    .command('create-user', 'Create a new user account')
    .action(async function (this: any) {
      try {
        // Prompt for user details
        const answers = await inquirer.prompt([
          {
            type: 'input',
            name: 'email',
            message: 'Email:',
            validate: (input: string) =>
              isValidEmail(input) || validationMessages.email,
          },
          {
            type: 'password',
            name: 'password',
            message: 'Password:',
            mask: '*',
            validate: (input: string) =>
              isValidPassword(input) || validationMessages.password,
          },
          {
            type: 'input',
            name: 'username',
            message: 'Username (optional):',
          },
        ]);

        const { prisma } = getContext();

        // Check if email already exists
        const existingUser = await prisma.user.findUnique({
          where: { email: answers.email },
        });

        if (existingUser) {
          output.error('User with this email already exists');
          return;
        }

        // Hash password
        const passwordHash = await bcrypt.hash(answers.password, 10);

        // Create user
        const spin = spinner('Creating user...');
        const user = await prisma.user.create({
          data: {
            email: answers.email,
            passwordHash,
            username: answers.username || null,
          },
        });
        spin.succeed('User created successfully!');

        // Display user details
        output.newline();
        output.section('User Details');
        console.log(`   ID:       ${user.id}`);
        console.log(`   Email:    ${user.email}`);
        console.log(`   Username: ${user.username || '(not set)'}`);
        console.log(`   Created:  ${formatDate(user.createdAt)}`);
        output.newline();
      } catch (error) {
        output.error(
          'Failed to create user: ' +
            (error instanceof Error ? error.message : error)
        );
      }
    });

  // List users command
  vorpal
    .command('list-users', 'List all users')
    .option('--limit <n>', 'Limit number of results')
    .option('--format <type>', 'Output format: table, json, csv')
    .action(async function (this: any, args: any) {
      try {
        const { prisma } = getContext();

        const limit = args.options.limit ? parseInt(args.options.limit) : undefined;
        const format = args.options.format || 'table';

        const users = await prisma.user.findMany({
          take: limit,
          orderBy: { createdAt: 'desc' },
        });

        if (users.length === 0) {
          output.warning('No users found');
          return;
        }

        if (format === 'json') {
          console.log(
            JSON.stringify(
              users.map(u => ({
                id: u.id,
                email: u.email,
                username: u.username,
                createdAt: u.createdAt,
              })),
              null,
              2
            )
          );
        } else if (format === 'csv') {
          console.log('ID,Email,Username,Created');
          users.forEach(u => {
            console.log(
              `${u.id},${u.email},${u.username || ''},${u.createdAt.toISOString()}`
            );
          });
        } else {
          // Table format (default)
          const table = createTable(['ID', 'Email', 'Username', 'Created']);
          users.forEach(u => {
            table.push([
              u.id.substring(0, 8) + '...',
              u.email,
              u.username || '(none)',
              formatDate(u.createdAt),
            ]);
          });
          console.log(table.toString());
          output.newline();
          output.info(`Total: ${users.length} user${users.length === 1 ? '' : 's'}`);
        }
      } catch (error) {
        output.error(
          'Failed to list users: ' +
            (error instanceof Error ? error.message : error)
        );
      }
    });

  // Delete user command
  vorpal
    .command('delete-user <identifier>', 'Delete a user by email or ID')
    .action(async function (this: any, args: any) {
      try {
        const { prisma } = getContext();
        const identifier = args.identifier;

        // Find user by email or ID
        const user = isValidUUID(identifier)
          ? await prisma.user.findUnique({ where: { id: identifier } })
          : await prisma.user.findUnique({ where: { email: identifier } });

        if (!user) {
          output.error('User not found');
          return;
        }

        // Get related data counts
        const [accordCount, recipeCount] = await Promise.all([
          prisma.accord.count({ where: { userId: user.id } }),
          prisma.recipe.count({ where: { userId: user.id } }),
        ]);

        // Show warning
        output.warning(
          'This will permanently delete the user and all their data!'
        );
        console.log(`   Email:    ${user.email}`);
        console.log(`   Username: ${user.username || '(none)'}`);
        console.log(`   Accords:  ${accordCount}`);
        console.log(`   Recipes:  ${recipeCount}`);
        output.newline();

        // Confirm deletion
        const confirm = await inquirer.prompt([
          {
            type: 'confirm',
            name: 'confirmed',
            message: 'Are you sure you want to delete this user?',
            default: false,
          },
        ]);

        if (!confirm.confirmed) {
          output.info('Deletion cancelled');
          return;
        }

        // Delete user (cascade will delete related data)
        const spin = spinner('Deleting user...');
        await prisma.user.delete({ where: { id: user.id } });
        spin.succeed('User deleted successfully!');
      } catch (error) {
        output.error(
          'Failed to delete user: ' +
            (error instanceof Error ? error.message : error)
        );
      }
    });

  // Reset password command
  vorpal
    .command('reset-password <identifier>', "Reset a user's password")
    .action(async function (this: any, args: any) {
      try {
        const { prisma } = getContext();
        const identifier = args.identifier;

        // Find user by email or ID
        const user = isValidUUID(identifier)
          ? await prisma.user.findUnique({ where: { id: identifier } })
          : await prisma.user.findUnique({ where: { email: identifier } });

        if (!user) {
          output.error('User not found');
          return;
        }

        console.log(
          `User: ${user.email} (${user.username || 'no username'})`
        );
        output.newline();

        // Prompt for new password
        const answers = await inquirer.prompt([
          {
            type: 'password',
            name: 'password',
            message: 'New password:',
            mask: '*',
            validate: (input: string) =>
              isValidPassword(input) || validationMessages.password,
          },
          {
            type: 'password',
            name: 'confirmPassword',
            message: 'Confirm password:',
            mask: '*',
            validate: (input: string, answers: any) =>
              input === answers.password || 'Passwords do not match',
          },
        ]);

        // Hash and update password
        const spin = spinner('Updating password...');
        const passwordHash = await bcrypt.hash(answers.password, 10);
        await prisma.user.update({
          where: { id: user.id },
          data: { passwordHash },
        });
        spin.succeed('Password reset successfully!');
      } catch (error) {
        output.error(
          'Failed to reset password: ' +
            (error instanceof Error ? error.message : error)
        );
      }
    });

  // Show user command
  vorpal
    .command('show-user <identifier>', 'Show detailed user information')
    .action(async function (this: any, args: any) {
      try {
        const { prisma } = getContext();
        const identifier = args.identifier;

        // Find user by email or ID
        const user = isValidUUID(identifier)
          ? await prisma.user.findUnique({ where: { id: identifier } })
          : await prisma.user.findUnique({ where: { email: identifier } });

        if (!user) {
          output.error('User not found');
          return;
        }

        // Get related data counts
        const [accordCount, recipeCount, refreshTokenCount] = await Promise.all([
          prisma.accord.count({ where: { userId: user.id } }),
          prisma.recipe.count({ where: { userId: user.id } }),
          prisma.refreshToken.count({ where: { userId: user.id } }),
        ]);

        // Display user details
        output.newline();
        output.section('👤 User Details');
        console.log(`ID:       ${user.id}`);
        console.log(`Email:    ${user.email}`);
        console.log(`Username: ${user.username || '(not set)'}`);
        console.log(`Created:  ${formatDate(user.createdAt)}`);
        console.log(`Updated:  ${formatDate(user.updatedAt)}`);

        output.section('📊 Statistics');
        console.log(`Accords:  ${accordCount}`);
        console.log(`Recipes:  ${recipeCount}`);
        console.log(`Active Refresh Tokens: ${refreshTokenCount}`);
        output.newline();
      } catch (error) {
        output.error(
          'Failed to show user: ' +
            (error instanceof Error ? error.message : error)
        );
      }
    });
}
