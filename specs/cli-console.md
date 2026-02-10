# CLI Console Specification

**Created**: February 10, 2026  
**Phase**: 13  
**Status**: Planning  
**Version**: 1.0

---

## Overview

Interactive command-line interface (REPL) for Scentora administration, providing direct access to CRUD operations without needing the web UI. Enables administrators to manage users, import/export data, and perform maintenance tasks efficiently.

### Purpose

Enable administrators and developers to:
- Create and manage users from the command line
- Import/export accords and recipes in bulk
- Perform database operations and maintenance
- Debug and test application functionality
- Seed development/test environments

### Key Design Principles

1. **User-Friendly Commands** - Simple, memorable command syntax
2. **Interactive Prompts** - Guide users through complex operations
3. **Rich UX** - Auto-complete, history, colored output, progress indicators
4. **Type-Safe** - TypeScript with full type checking
5. **Reusable** - Leverage existing services and repositories
6. **Extensible** - Easy to add new commands and features

---

## Technology Stack

### Core Framework
- **Vorpal.js** - Interactive CLI framework with REPL support
  - Command parsing and routing
  - Auto-complete and command history
  - Interactive prompts (inquirer.js integration)
  - TypeScript support

### Supporting Libraries
- **chalk** - Terminal string styling (colors, bold, etc.)
- **cli-table3** - Beautiful ASCII tables
- **ora** - Elegant terminal spinners
- **tsx** - TypeScript execution (consistent with backend)
- **inquirer** - Interactive prompts (used by Vorpal)

### Architecture
```
backend/cli/
├── index.ts              # Entry point, bootstrap REPL
├── context.ts            # Application context (DB, services, config)
├── commands/
│   ├── user.ts          # User management commands
│   ├── accord.ts        # Accord management commands (future)
│   ├── recipe.ts        # Recipe management commands (future)
│   └── db.ts            # Database operations (future)
├── utils/
│   ├── output.ts        # Formatted output helpers
│   ├── validation.ts    # Input validation
│   └── prompts.ts       # Reusable prompt definitions
└── types/
    └── context.ts       # Context type definitions
```

---

## Commands

### Phase 13.1: User Management (Initial Release)

#### `create-user`
Create a new user account interactively.

**Usage**:
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

**Validation**:
- Email format validation
- Password strength requirements (8+ chars)
- Username uniqueness check
- Email uniqueness check

**Implementation**:
- Prompts for email, password, username
- Hashes password with bcrypt
- Creates user via `userService.create()`
- Returns formatted success/error message

---

#### `list-users`
Display all users in a formatted table.

**Usage**:
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

**Options**:
- `--limit <n>` - Limit results (default: all)
- `--format <type>` - Output format: table (default), json, csv

**Implementation**:
- Fetches all users via Prisma
- Formats as ASCII table with cli-table3
- Supports alternate formats (JSON, CSV)

---

#### `delete-user <email|id>`
Delete a user by email or ID with confirmation.

**Usage**:
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

**Validation**:
- User must exist
- Requires confirmation prompt
- Shows data that will be deleted (accords, recipes count)

**Implementation**:
- Finds user by email or ID
- Displays user details and related data counts
- Prompts for confirmation
- Deletes user (cascade deletes accords, recipes per FK constraints)
- Shows spinner during deletion

---

#### `reset-password <email|id>`
Reset a user's password interactively.

**Usage**:
```bash
scentora> reset-password admin@example.com
User: admin@example.com (admin)
? New password: ********
? Confirm password: ********
✅ Password reset successfully!
```

**Validation**:
- User must exist
- Password strength requirements
- Passwords must match

**Implementation**:
- Finds user by email or ID
- Prompts for new password (hidden input)
- Prompts for confirmation
- Hashes password with bcrypt
- Updates user via Prisma
- Optionally invalidates refresh tokens

---

#### `show-user <email|id>`
Display detailed information about a specific user.

**Usage**:
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

**Implementation**:
- Finds user by email or ID
- Fetches related data counts (accords, recipes, tokens)
- Formats as structured output with sections

---

### Phase 13.2: Accord Management (Future)

#### `import-accords <file>`
Import accords from JSON file.

**Usage**:
```bash
scentora> import-accords ./accords-backup.json
📂 Reading file: ./accords-backup.json
Found 50 accords
? Assign to user (email): admin@example.com
📥 Importing accords... ████████████████████ 50/50
✅ Import complete!
   Created: 50
   Skipped: 0
   Errors: 0
```

---

#### `export-accords <email|id> [file]`
Export user's accords to JSON file.

**Usage**:
```bash
scentora> export-accords admin@example.com ./backup.json
📤 Exporting 15 accords... done
✅ Export complete: ./backup.json (23.4 KB)
```

---

#### `list-accords <email|id>`
List all accords for a user.

---

### Phase 13.3: Recipe Management (Future)

#### `import-recipes <file>`
Import recipes from JSON file.

#### `export-recipes <email|id> [file]`
Export user's recipes to JSON file.

#### `clone-recipe <id> <target-user>`
Clone a recipe to another user.

---

### Phase 13.4: Database Operations (Future)

#### `db:migrate`
Run pending database migrations.

#### `db:seed`
Seed database with test data.

#### `db:backup [file]`
Create database backup.

#### `db:restore <file>`
Restore database from backup.

---

## Application Context

### Context Loading
On startup, the CLI loads application context including:
- Database connection (Prisma client)
- Environment configuration
- Service instances (UserService, AccordService, etc.)
- Logger instance

### Context Interface
```typescript
interface CLIContext {
  prisma: PrismaClient;
  services: {
    user: UserService;
    accord: AccordService;
    recipe: RecipeService;
  };
  config: {
    databaseUrl: string;
    environment: string;
  };
  logger: Logger;
}
```

### Error Handling
- Database connection errors → show friendly message, exit gracefully
- Command errors → show error message, don't crash REPL
- Validation errors → show field-specific errors
- Unexpected errors → log stack trace, show user-friendly message

---

## User Experience Features

### Auto-Complete
- Command names (e.g., typing `cre<TAB>` completes to `create-user`)
- Command options (e.g., typing `--f<TAB>` completes to `--format`)
- Email addresses (from database when relevant)

### Command History
- Up/down arrows navigate through command history
- History persisted across sessions (`.scentora_history`)
- Search history with Ctrl+R

### Colored Output
- ✅ Success messages → green
- ❌ Error messages → red
- ⚠️ Warnings → yellow
- 📊 Info messages → blue
- User input → cyan
- Tables → mix of colors for headers/data

### Progress Indicators
- Spinners for long-running operations (ora)
- Progress bars for batch operations (cli-progress)
- ETA for large imports/exports

### Formatted Tables
- ASCII tables with borders (cli-table3)
- Column alignment (left for text, right for numbers)
- Truncation for long values with ellipsis
- Header styling (bold, colored)

---

## Entry Point & Execution

### Package.json Scripts
```json
{
  "scripts": {
    "console": "tsx cli/index.ts",
    "cli": "tsx cli/index.ts"
  }
}
```

### Usage
```bash
# From backend directory
npm run console

# Or with arguments
npm run console -- --help
npm run console -- --version
```

### Startup Sequence
1. Load environment variables from `.env`
2. Connect to database (Prisma)
3. Initialize services
4. Load command modules
5. Display welcome banner
6. Start Vorpal REPL

### Welcome Banner
```
🌸 Scentora CLI Console (v1.0.0)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Environment: development
Database: localhost:5432/scentora
Connected: ✅

Type 'help' for available commands
Type 'exit' to quit
```

---

## Testing Strategy

### Unit Tests
- Test individual command logic (without REPL)
- Mock Prisma client and services
- Test validation functions
- Test output formatters

### Integration Tests
- Test commands end-to-end with test database
- Verify database operations
- Test error handling
- Test prompts and confirmations

### Manual Testing
- Run CLI and verify all commands work
- Test auto-complete and history
- Verify colored output renders correctly
- Test edge cases (empty database, invalid input)

---

## Security Considerations

### Authentication
- No authentication required (assumes CLI user has server access)
- CLI should only be accessible to server administrators
- Document that CLI is for admin use only

### Password Handling
- Passwords never logged or displayed
- Use hidden input for password prompts
- Hash passwords before storing

### Input Validation
- Sanitize all user input
- Validate email formats
- Prevent SQL injection (Prisma handles this)
- Validate file paths for import/export

### Audit Logging
- Log all CLI commands executed (optional feature)
- Include timestamp, user, command, result
- Store in separate audit log file

---

## Development Phases

### Phase 13.1: Foundation & User Management ✅
**Goal**: Working CLI with user management commands

**Tasks**:
1. Setup CLI structure and entry point
2. Configure Vorpal with TypeScript
3. Implement application context loading
4. Create user management commands:
   - `create-user`
   - `list-users`
   - `delete-user`
   - `reset-password`
   - `show-user`
5. Add output utilities (colors, tables, spinners)
6. Write tests for user commands
7. Document CLI usage

**Deliverables**:
- `backend/cli/` directory with working CLI
- 5 user management commands
- Tests for all commands
- Documentation in `docs/CLI_CONSOLE.md`
- Updated `package.json` with console script

---

### Phase 13.2: Accord Management (Future)
**Goal**: Import/export and manage accords via CLI

**Commands**:
- `import-accords`
- `export-accords`
- `list-accords`
- `delete-accord`
- `bulk-update-accords`

---

### Phase 13.3: Recipe Management (Future)
**Goal**: Import/export and manage recipes via CLI

**Commands**:
- `import-recipes`
- `export-recipes`
- `list-recipes`
- `clone-recipe`
- `delete-recipe`

---

### Phase 13.4: Database Operations (Future)
**Goal**: Database maintenance and seeding

**Commands**:
- `db:migrate`
- `db:seed`
- `db:backup`
- `db:restore`
- `db:reset`

---

### Phase 13.5: Advanced Features (Future)
**Goal**: Enhanced UX and power features

**Features**:
- Command aliases (shortcuts)
- Batch operations (process multiple items)
- Query DSL for filtering
- Export to multiple formats (CSV, YAML, etc.)
- Scripting mode (non-interactive)
- Audit logging
- Configuration file (.scentorarc)

---

## Success Criteria

### Phase 13.1 Success Criteria
- ✅ CLI starts without errors
- ✅ All 5 user commands functional
- ✅ Auto-complete works for commands
- ✅ Command history persists across sessions
- ✅ Colored output renders correctly
- ✅ Tables display properly formatted
- ✅ Password prompts hide input
- ✅ Confirmation prompts work correctly
- ✅ Error messages are clear and helpful
- ✅ Tests pass for all commands (80%+ coverage)
- ✅ Documentation complete and accurate

### Overall Project Success
- Intuitive, user-friendly command interface
- Fast and responsive (commands complete in <2s)
- Robust error handling with clear messages
- Comprehensive test coverage
- Well-documented with examples
- Easy to extend with new commands

---

## Future Enhancements

### Potential Features
- **Interactive Mode**: Menu-driven interface for common tasks
- **Scripting Support**: Run CLI commands from script files
- **Remote Access**: Connect to remote databases
- **Multi-tenancy**: Switch between different environments
- **Plugin System**: Allow third-party command extensions
- **GraphQL Console**: Interactive GraphQL query interface
- **Performance Profiling**: Built-in performance monitoring
- **Data Migration Tools**: Convert between schema versions
- **Batch Processing**: Process large datasets efficiently
- **Scheduled Tasks**: Run commands on a schedule (cron-like)

---

## Resources & References

### Libraries
- **Vorpal**: https://github.com/dthree/vorpal
- **chalk**: https://github.com/chalk/chalk
- **cli-table3**: https://github.com/cli-table/cli-table3
- **ora**: https://github.com/sindresorhus/ora
- **inquirer**: https://github.com/SBoudrias/Inquirer.js

### Inspirations
- Rails Console (Ruby on Rails)
- Laravel Tinker (Laravel/PHP)
- Django Shell (Django/Python)
- Prisma Studio CLI
- MongoDB Shell

---

**Phase**: 13  
**Document**: specs/cli-console.md  
**Version**: 1.0  
**Last Updated**: February 10, 2026
