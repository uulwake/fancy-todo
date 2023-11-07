import { NextFunction, Request, Response } from "express";
import jwt from "jsonwebtoken";

import { CustomError } from "../libs/custom-error";

export const authenticateJwt = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const jwtToken = req.headers.jwt_token;
  if (!jwtToken) {
    next(new CustomError({ message: "Missing JWT Token", status: 401 }));
    return;
  }

  try {
    const decoded = jwt.verify(jwtToken, process.env.JWT_SECRET) as {
      id: number;
    };

    req.user_id = decoded.id;
    next();
  } catch (err) {
    const message = err instanceof Error ? err.message : "Invalid JWT Token";

    next(
      new CustomError({
        message,
        status: 401,
      })
    );
  }
};
