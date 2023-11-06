declare global {
  namespace NodeJS {
    interface ProcessEnv {
      DATABASE_URL: string;
      PORT: string;
      SALT: string;
      JWT_SECRET: string;
      JWT_EXPIRED: string;
      TZ: string;
    }
  }
}

export {};
