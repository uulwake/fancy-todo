import { NextFunction, Request, Response } from "express";
import { validationResult } from "express-validator";
import { CustomError } from "../libs/custom-error";

export const reqValidationResult = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const err = validationResult(req).array();
  if (err && err.length) {
    next(
      new CustomError({
        message: err.map((el) => el.msg).join(", "),
        status: 400,
      })
    );
  } else {
    next();
  }
};
