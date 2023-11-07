export type TaskModel = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  status: string;
  order: number;
  created_at: Date;
  updated_at: Date;
};

export type TaskModelField =
  | "id"
  | "user_id"
  | "title"
  | "description"
  | "status"
  | "order"
  | "created_at"
  | "updated_at";
