import nano from 'nano';
import { config } from './index';

const couchUrl = `${config.couchdb.url.replace('://', `://${config.couchdb.user}:${config.couchdb.password}@`)}`;
const client = nano(couchUrl);

export const db = client.db.use(config.couchdb.database);

export async function initDatabase() {
  try {
    const dbList = await client.db.list();
    
    if (!dbList.includes(config.couchdb.database)) {
      await client.db.create(config.couchdb.database);
      console.log(`✓ Created database: ${config.couchdb.database}`);
    } else {
      console.log(`✓ Database exists: ${config.couchdb.database}`);
    }

    // Create indexes for efficient querying
    await db.createIndex({
      index: {
        fields: ['type', 'userId', 'createdAt'],
      },
      name: 'type-userId-createdAt-index',
    });

    await db.createIndex({
      index: {
        fields: ['type', 'userId', 'perfumeId', 'date'],
      },
      name: 'journal-by-user-perfume-index',
    });

    await db.createIndex({
      index: {
        fields: ['type', 'email'],
      },
      name: 'user-email-index',
    });

    await db.createIndex({
      index: {
        fields: ['type', 'username'],
      },
      name: 'user-username-index',
    });

    await db.createIndex({
      index: {
        fields: ['type', 'code'],
      },
      name: 'invitation-code-index',
    });

    await db.createIndex({
      index: {
        fields: ['type', 'createdBy', 'createdAt'],
      },
      name: 'invitation-by-creator-index',
    });

    console.log('✓ Database indexes created');
    
  } catch (error) {
    console.error('Database initialization error:', error);
    throw error;
  }
}
