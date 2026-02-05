import type Koa from 'koa';

export interface AuthState {
  userId: string;
  email: string;
}

export type AuthContext = Koa.ParameterizedContext<AuthState>;
