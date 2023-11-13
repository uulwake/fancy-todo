import { NextFunction, Request, Response } from "express";
import { v4 as uuidV4 } from "uuid";

export default () => {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      req.headers.request_id = uuidV4();
      next();
    } catch (err) {
      next(err);
    }
  };
};
