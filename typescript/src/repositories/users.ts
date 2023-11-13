import { IUserRepo } from "./interfaces";
import { UserModel, UserModelField } from "../models/user";
import { DBType } from "../databases";
import { Context } from "../types/context";

export class UserRepo implements IUserRepo {
  private db: DBType;

  constructor(db: DBType) {
    this.db = db;
  }

  async createUser(ctx: Context, data: Omit<UserModel, "id">): Promise<number> {
    const res = await this.db
      .pg<UserModel>("users")
      .insert(data)
      .returning("id");

    return res[0].id;
  }

  async getUserDetail(ctx: Context, opt: {
    id?: number;
    email?: string;
    cols?: UserModelField[];
  }): Promise<Partial<UserModel>> {
    const cols = opt.cols ?? ["id"];
    const query = this.db.pg("users").select(cols).first();

    if (opt.id) {
      query.where("id", opt.id);
    }

    if (opt.email) {
      query.where("email", opt.email);
    }

    return query as Partial<UserModel>;
  }
}
