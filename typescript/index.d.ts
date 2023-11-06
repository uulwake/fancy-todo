import { IncomingHttpHeaders } from "http";

declare global {
  namespace Express {
    export interface Request {
      user_id?: number;
    }
  }
}

declare module "http" {
  interface IncomingHttpHeaders {
    jwt_token?: string;
  }
}

export {};
