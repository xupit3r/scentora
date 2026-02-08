import { parseArgs } from 'node:util';
import { createInterface } from 'node:readline';
import bcrypt from 'bcryptjs';
import { prisma } from './prisma.js';
import { createInvitation } from './services/invitation.service.js';

function printUsage() {
  console.log(`Usage: node dist/cli.js <command> [options]

Commands:
  bootstrap            Create the first user and an invitation code
  create-invitation    Create an invitation code for an existing user

Run '<command> --help' for command-specific options.`);
}

function printBootstrapHelp() {
  console.log(`Usage: node dist/cli.js bootstrap --email <email> --username <username>

Creates the first user account and a 30-day invitation code.
Refuses to run if any users already exist (safety check).

Password is prompted interactively.

Options:
  --email     User email address (required)
  --username  Username (required)
  --help      Show this help`);
}

function printCreateInvitationHelp() {
  console.log(`Usage: node dist/cli.js create-invitation --email <creator-email> [options]

Creates an invitation code using an existing user as the creator.

Options:
  --email      Creator's email address to look up user (required)
  --for-email  Restrict invitation to a specific email
  --days       Expiration in days (default: 7)
  --help       Show this help`);
}

async function promptPassword(prompt: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const rl = createInterface({
      input: process.stdin,
      output: process.stderr,
      terminal: process.stdin.isTTY ?? false,
    });

    if (process.stdin.isTTY) {
      // Hide input by using raw mode
      process.stderr.write(prompt);
      let password = '';
      process.stdin.setRawMode(true);
      process.stdin.resume();
      process.stdin.setEncoding('utf8');

      const onData = (ch: string) => {
        if (ch === '\n' || ch === '\r' || ch === '\u0004') {
          process.stdin.setRawMode(false);
          process.stdin.pause();
          process.stdin.removeListener('data', onData);
          process.stderr.write('\n');
          rl.close();
          resolve(password);
        } else if (ch === '\u0003') {
          // Ctrl+C
          process.stdin.setRawMode(false);
          rl.close();
          reject(new Error('Aborted'));
        } else if (ch === '\u007f' || ch === '\b') {
          // Backspace
          password = password.slice(0, -1);
        } else {
          password += ch;
        }
      };
      process.stdin.on('data', onData);
    } else {
      // Non-interactive: read a line from stdin
      rl.question(prompt, (answer) => {
        rl.close();
        resolve(answer);
      });
    }
  });
}

async function bootstrapCommand(args: string[]) {
  const { values } = parseArgs({
    args,
    options: {
      email: { type: 'string' },
      username: { type: 'string' },
      help: { type: 'boolean', short: 'h' },
    },
    strict: true,
  });

  if (values.help) {
    printBootstrapHelp();
    return;
  }

  if (!values.email || !values.username) {
    console.error('Error: --email and --username are required.');
    printBootstrapHelp();
    process.exitCode = 1;
    return;
  }

  // Safety check: refuse if users already exist
  const userCount = await prisma.user.count();
  if (userCount > 0) {
    console.error(`Error: ${userCount} user(s) already exist. Bootstrap is only for first-time setup.`);
    process.exitCode = 1;
    return;
  }

  const password = await promptPassword('Password: ');
  if (!password) {
    console.error('Error: password is required.');
    process.exitCode = 1;
    return;
  }

  const passwordHash = await bcrypt.hash(password, 10);
  const user = await prisma.user.create({
    data: {
      email: values.email,
      username: values.username,
      passwordHash,
    },
  });

  const invitation = await createInvitation(user.id, undefined, 30);

  console.log('');
  console.log('=== Bootstrap Complete ===');
  console.log(`  User:       ${user.email} (${user.username})`);
  console.log(`  Invitation: ${invitation.code}`);
  console.log(`  Expires:    30 days`);
  console.log('');
  console.log('Log in with the user above, or register a new account');
  console.log('using the invitation code.');
}

async function createInvitationCommand(args: string[]) {
  const { values } = parseArgs({
    args,
    options: {
      email: { type: 'string' },
      'for-email': { type: 'string' },
      days: { type: 'string', short: 'd' },
      help: { type: 'boolean', short: 'h' },
    },
    strict: true,
  });

  if (values.help) {
    printCreateInvitationHelp();
    return;
  }

  if (!values.email) {
    console.error('Error: --email is required (creator\'s email).');
    printCreateInvitationHelp();
    process.exitCode = 1;
    return;
  }

  const days = values.days ? parseInt(values.days, 10) : 7;
  if (isNaN(days) || days < 1) {
    console.error('Error: --days must be a positive integer.');
    process.exitCode = 1;
    return;
  }

  const user = await prisma.user.findUnique({ where: { email: values.email } });
  if (!user) {
    console.error(`Error: no user found with email "${values.email}".`);
    process.exitCode = 1;
    return;
  }

  const forEmail = values['for-email'];
  const invitation = await createInvitation(user.id, forEmail, days);

  console.log('');
  console.log('Invitation created!');
  console.log(`  Code:    ${invitation.code}`);
  if (forEmail) console.log(`  For:     ${forEmail}`);
  console.log(`  Expires: ${invitation.expiresAt}`);
}

async function main() {
  const args = process.argv.slice(2);
  const command = args[0];

  if (!command || command === '--help' || command === '-h') {
    printUsage();
    return;
  }

  const commandArgs = args.slice(1);

  switch (command) {
    case 'bootstrap':
      await bootstrapCommand(commandArgs);
      break;
    case 'create-invitation':
      await createInvitationCommand(commandArgs);
      break;
    default:
      console.error(`Unknown command: ${command}`);
      printUsage();
      process.exitCode = 1;
  }
}

main()
  .catch((err) => {
    console.error(err);
    process.exitCode = 1;
  })
  .finally(() => prisma.$disconnect());
