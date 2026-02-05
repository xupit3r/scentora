import type Koa from 'koa';
import authRoutes from './auth.routes.js';
import accordRoutes from './accord.routes.js';
import tagRoutes from './tag.routes.js';
import statsRoutes from './stats.routes.js';
import exportRoutes from './export.routes.js';
import invitationRoutes from './invitation.routes.js';
import recipeRoutes from './recipe.routes.js';
import collectionRoutes from './collection.routes.js';

export function mountRoutes(app: Koa) {
  app.use(authRoutes.routes()).use(authRoutes.allowedMethods());
  app.use(accordRoutes.routes()).use(accordRoutes.allowedMethods());
  app.use(tagRoutes.routes()).use(tagRoutes.allowedMethods());
  app.use(statsRoutes.routes()).use(statsRoutes.allowedMethods());
  app.use(exportRoutes.routes()).use(exportRoutes.allowedMethods());
  app.use(invitationRoutes.routes()).use(invitationRoutes.allowedMethods());
  app.use(recipeRoutes.routes()).use(recipeRoutes.allowedMethods());
  app.use(collectionRoutes.routes()).use(collectionRoutes.allowedMethods());
}
