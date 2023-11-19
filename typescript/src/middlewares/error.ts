import { NextFunction, Request, Response } from "express";
import { PG_ERROR_CODE } from "../constants";

export default () => {
  return (
    err: { message: string; status?: number; [key: string]: any },
    req: Request,
    res: Response,
    next: NextFunction
  ) => {
    let message = err.message || "Server Error";
    let status = err.status || 500;

    res.status(status).json({ message });
  };
};
