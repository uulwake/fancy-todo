import { NextFunction, Request, Response } from "express";
import { PG_ERROR_CODE } from "../constants";

export default (
  err: { message: string; status?: number; [key: string]: any },
  req: Request,
  res: Response,
  next: NextFunction
) => {
  let message = err.message || "Server Error";
  let status = err.status || 500;

  if (err.schema && err.table) {
    if (err.code === PG_ERROR_CODE.CONSTRAINT.UNIQUE) {
      const matches = err.detail.matchAll(/\(([\w@\.]+)\)/gim);
      message = `${matches.next().value[1]} ${
        matches.next().value[1]
      } already exists`;
    } else {
      message = err.detail;
    }
  }

  res.status(status).json({ message });
};
