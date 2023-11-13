import { Request } from "express";
import { Context } from "../types/context";

export const createContext = (req: Request): Context => {
  return {
    ipAddr: req.clientIp,
    requestId: req.headers.request_id,
    userId: req.user_id,
  };
};
