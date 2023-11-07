export type UserModel = {
  id: number;
  name: string;
  email: string;
  password: string;
  created_at: Date;
  updated_at: Date;
};

export type UserModelField =
  | "id"
  | "name"
  | "email"
  | "password"
  | "created_at"
  | "updated_at";
