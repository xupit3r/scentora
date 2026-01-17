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

    // Create design documents/views here in the future
    
  } catch (error) {
    console.error('Database initialization error:', error);
    throw error;
  }
}
