import { NextFunction, Request, Response, Router } from "express";

import { UserModel } from "../models/user";
import { ServiceType } from "../services/types";
import { BaseHandler } from "./base";
import { UserHandlerValidatorType } from "./validators/types";
import { createContext } from "../libs/create-context";

export class UserHandler extends BaseHandler {
  private service: ServiceType;
  private validator: UserHandlerValidatorType;

  constructor(service: ServiceType, validator: UserHandlerValidatorType) {
    super(Router());

    this.service = service;
    this.validator = validator;

    this.router.post(
      "/register",
      this.validator.register,
      this.register.bind(this)
    );

    this.router.post("/login", this.validator.login, this.login.bind(this));
  }

  async register(
    req: Request<{}, {}, Pick<UserModel, "name" | "email" | "password">>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const { id, jwt_token } = await this.service.userService.register(
        createContext(req),
        req.body
      );
      res.json({ data: { user: { id }, jwt_token } });
    } catch (err) {
      next(err);
    }
  }

  async login(
    req: Request<{}, {}, Pick<UserModel, "email" | "password">>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const { id, jwt_token } = await this.service.userService.login(
        createContext(req),
        req.body
      );
      res.json({ data: { user: { id }, jwt_token } });
    } catch (err) {
      next(err);
    }
  }
}
