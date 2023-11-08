declare global {
  namespace NodeJS {
    interface ProcessEnv {
      NODE_ENV: string;
      DATABASE_URL: string;
      ES_URL: string;
      PORT: string;
      SALT: string;
      JWT_SECRET: string;
      JWT_EXPIRED: string;
      TZ: string;
    }
  }
}

export {};
