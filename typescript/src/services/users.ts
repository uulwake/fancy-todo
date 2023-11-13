import jwt from "jsonwebtoken";
import bcrypt from "bcrypt";

import { UserModel } from "../models/user";
import { RepositoryType } from "../repositories/types";
import { CustomError } from "../libs/custom-error";
import { IUserService } from "./interfaces";
import { Context } from "../types/context";

export class UserService implements IUserService {
  private repo;
  constructor(repo: RepositoryType) {
    this.repo = repo;
  }

  private createUserToken(ctx: Context, data: { id: number; email: string }): string {
    return jwt.sign(data, process.env.JWT_SECRET, {
      expiresIn: process.env.JWT_EXPIRED,
    });
  }

  async register(
    ctx: Context,
    body: Pick<UserModel, "name" | "email" | "password">
  ): Promise<{ id: number; jwt_token: string }> {
    const now = new Date();

    const userData: Omit<UserModel, "id"> = {
      name: body.name,
      email: body.email,
      password: await bcrypt.hash(body.password, Number(process.env.SALT)),
      created_at: now,
      updated_at: now,
    };

    const id = await this.repo.userRepo.createUser(ctx, userData);
    const jwt_token = this.createUserToken(ctx, { id, email: userData.email });

    return { id, jwt_token };
  }

  async login(
    ctx: Context,
    body: Pick<UserModel, "email" | "password">
  ): Promise<{ id: number; jwt_token: string }> {
    const user = await this.repo.userRepo.getUserDetail(ctx, {
      email: body.email,
      cols: ["id", "email", "password"],
    });

    if (!user || !user.id || !user.email || !user.password) {
      throw new CustomError({ message: "invalid email/password", status: 404 });
    }

    const isValid = await bcrypt.compare(body.password, user.password);
    if (!isValid) {
      throw new CustomError({ message: "invalid email/password", status: 404 });
    }

    const jwt_token = this.createUserToken(ctx, { id: user.id, email: user.email });
    return { id: user.id, jwt_token };
  }
}
