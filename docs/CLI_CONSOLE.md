# Scentora CLI Console

Interactive command-line interface for Scentora administration and data management.

**Version**: 1.0  
**Last Updated**: February 10, 2026

---

## Overview

The Scentora CLI Console provides a powerful, interactive command-line interface for managing users, accords, recipes, and performing database operations without needing the web UI.

### Features

- 🎯 **User-Friendly Commands** - Simple, intuitive command syntax
- 🎨 **Rich Terminal UI** - Colored output, tables, spinners, progress bars
- ⌨️ **Auto-Complete** - Tab completion for commands and options
- 📜 **Command History** - Navigate previous commands with arrow keys
- 🔒 **Secure** - Password hashing, input validation, confirmation prompts
- 🧪 **Well-Tested** - Comprehensive test coverage for reliability

---

## Installation

The CLI console is included with the Scentora backend. No additional installation required.

### Prerequisites

- Node.js 18+
- Scentora backend installed
- Database configured and running

---

## Getting Started

### Starting the Console

From the `backend` directory:

```bash
npm run console
```

You'll see the welcome banner:

```
🌸 Scentora CLI Console (v1.0.0)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Environment: development
Database: localhost:5432/scentora
Connected: ✅

Type 'help' for available commands
Type 'exit' to quit

scentora>
```

### Getting Help

```bash
scentora> help

  Commands:

    help [command...]       Provides help for a given command.
    exit                    Exits application.
    create-user            Create a new user account
    list-users             List all users
    delete-user <id>       Delete a user by email or ID
    reset-password <id>    Reset a user's password
    show-user <id>         Show detailed user information
```

For help with a specific command:

```bash
scentora> help create-user
```

---

## Commands

### User Management

#### `create-user`

Create a new user account with interactive prompts.

**Usage**:
```bash
scentora> create-user
```

**Interactive Prompts**:
1. Email address (validated format)
2. Password (hidden input, min 8 characters)
3. Username (optional)

**Example**:
```bash
scentora> create-user
? Email: admin@example.com
? Password: ********
? Username: admin
✅ User created successfully!
   ID: 550e8400-e29b-41d4-a716-446655440000
   Email: admin@example.com
   Username: admin
   Created: 2026-02-10T18:00:00Z
```

**Validation Rules**:
- Email must be valid format
- Email must be unique
- Password must be at least 8 characters
- Username must be unique (if provided)

---

#### `list-users`

Display all users in a formatted table.

**Usage**:
```bash
scentora> list-users [options]
```

**Options**:
- `--limit <n>` - Limit number of results
- `--format <type>` - Output format: `table` (default), `json`, `csv`

**Example**:
```bash
scentora> list-users
┌──────────────────────────────────────┬────────────────────┬──────────┬─────────────────────┐
│ ID                                   │ Email              │ Username │ Created             │
├──────────────────────────────────────┼────────────────────┼──────────┼─────────────────────┤
│ 550e8400-e29b-41d4-a716-446655440000 │ admin@example.com  │ admin    │ 2026-02-10 18:00:00 │
│ 6ba7b810-9dad-11d1-80b4-00c04fd430c8 │ user@example.com   │ user     │ 2026-02-09 14:30:00 │
└──────────────────────────────────────┴────────────────────┴──────────┴─────────────────────┘

Total: 2 users
```

**JSON Format**:
```bash
scentora> list-users --format json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@example.com",
    "username": "admin",
    "createdAt": "2026-02-10T18:00:00Z"
  }
]
```

---

#### `delete-user <email|id>`

Delete a user account with confirmation prompt.

**Usage**:
```bash
scentora> delete-user <email|id>
```

**Arguments**:
- `<email|id>` - User's email address or ID

**Example**:
```bash
scentora> delete-user admin@example.com
⚠️  Warning: This will permanently delete the user and all their data!
   Email: admin@example.com
   Username: admin
   Accords: 15
   Recipes: 8
? Are you sure you want to delete this user? (y/N) y
🗑️  Deleting user...
✅ User deleted successfully!
```

**Data Deletion**:
When a user is deleted, the following data is also deleted (cascade):
- All accords owned by the user
- All recipes owned by the user
- All recipe versions
- All recipe ingredients
- All tags created by the user
- All refresh tokens

---

#### `reset-password <email|id>`

Reset a user's password interactively.

**Usage**:
```bash
scentora> reset-password <email|id>
```

**Arguments**:
- `<email|id>` - User's email address or ID

**Example**:
```bash
scentora> reset-password admin@example.com
User: admin@example.com (admin)
? New password: ********
? Confirm password: ********
✅ Password reset successfully!
```

**Security**:
- Password input is hidden
- Password must meet strength requirements (8+ characters)
- Passwords must match
- Password is hashed with bcrypt before storage

---

#### `show-user <email|id>`

Display detailed information about a user.

**Usage**:
```bash
scentora> show-user <email|id>
```

**Arguments**:
- `<email|id>` - User's email address or ID

**Example**:
```bash
scentora> show-user admin@example.com

👤 User Details
────────────────────────────────────────
ID:       550e8400-e29b-41d4-a716-446655440000
Email:    admin@example.com
Username: admin
Created:  2026-02-10T18:00:00Z
Updated:  2026-02-10T18:00:00Z

📊 Statistics
────────────────────────────────────────
Accords:  15
Recipes:  8
Active Refresh Tokens: 2
```

---

## Terminal Features

### Auto-Complete

Press `Tab` to auto-complete commands and options:

```bash
scentora> cre<TAB>
scentora> create-user
```

### Command History

Navigate through previous commands:
- `↑` (Up Arrow) - Previous command
- `↓` (Down Arrow) - Next command
- `Ctrl+R` - Search command history

History is persisted across sessions in `.scentora_history`.

### Colored Output

- ✅ Success messages → Green
- ❌ Error messages → Red
- ⚠️ Warnings → Yellow
- 📊 Info messages → Blue
- User input → Cyan

### Progress Indicators

Long-running operations show progress:
- Spinners for operations (e.g., "Deleting user...")
- Progress bars for batch operations (e.g., "Importing 50/100...")

---

## Common Tasks

### Creating Your First User

```bash
scentora> create-user
? Email: myemail@example.com
? Password: ********
? Username: myusername
✅ User created successfully!
```

### Listing All Users

```bash
scentora> list-users
# See table of all users
```

### Resetting a Forgotten Password

```bash
scentora> reset-password user@example.com
? New password: ********
? Confirm password: ********
✅ Password reset successfully!
```

### Checking User Details

```bash
scentora> show-user user@example.com
# See detailed user information
```

---

## Tips & Best Practices

### 1. Use Tab Completion
Save time by using `Tab` to auto-complete commands and options.

### 2. Review Before Deleting
Always review the confirmation prompt when deleting users - this operation cannot be undone.

### 3. Use JSON Format for Scripting
Export data in JSON format for processing with other tools:
```bash
scentora> list-users --format json > users.json
```

### 4. Check Help for Any Command
```bash
scentora> help <command-name>
```

### 5. Exit Gracefully
Always use `exit` or `Ctrl+C` to close the console properly.

---

## Troubleshooting

### Console Won't Start

**Problem**: Error connecting to database

**Solution**:
1. Verify database is running: `docker compose ps`
2. Check DATABASE_URL in `.env`
3. Verify credentials are correct

---

### Command Not Found

**Problem**: "Command not found" error

**Solution**:
1. Type `help` to see available commands
2. Check spelling
3. Use `Tab` completion

---

### Cannot Delete User

**Problem**: "User has associated data" error

**Solution**: This is expected - delete operations cascade to related data. Review the confirmation prompt carefully.

---

### Password Requirements Not Met

**Problem**: "Password must be at least 8 characters"

**Solution**: Use a stronger password with at least 8 characters.

---

## Security Considerations

### Access Control
- The CLI console should only be accessible to server administrators
- Requires direct access to the server where Scentora is running
- No network-accessible authentication required

### Password Security
- Passwords are never displayed or logged
- All passwords are hashed with bcrypt before storage
- Hidden input for password prompts

### Audit Trail
All CLI operations should be logged for audit purposes (future feature).

---

## Development & Extension

### Adding New Commands

Commands are located in `backend/cli/commands/`. To add a new command:

1. Create a new file in `backend/cli/commands/`
2. Export a Vorpal command definition
3. Import in `backend/cli/index.ts`

See existing commands for examples.

### Running Tests

```bash
cd backend
npm test -- cli/
```

---

## Roadmap

### Upcoming Features (Phase 13.2+)

- **Accord Management**: Import/export accords, bulk operations
- **Recipe Management**: Import/export recipes, clone recipes
- **Database Operations**: Migrations, seeding, backups
- **Scripting Mode**: Run commands from script files
- **Audit Logging**: Track all CLI operations
- **Remote Access**: Connect to remote databases

---

## Support

For issues, questions, or feature requests:
- GitHub Issues: https://github.com/xupit3r/scentora/issues
- Documentation: https://github.com/xupit3r/scentora/docs

---

**Document**: docs/CLI_CONSOLE.md  
**Version**: 1.0  
**Last Updated**: February 10, 2026
