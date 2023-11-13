import { IncomingHttpHeaders } from "http";

declare global {
  namespace Express {
    export interface Request {
      user_id?: number;
      requestId: string;
      clientIp: string;
    }
  }
}

declare module "http" {
  interface IncomingHttpHeaders {
    jwt_token?: string;
    request_id: string;
  }
}

export {};
