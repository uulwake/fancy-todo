export type TaskTagModel = {
  task_id: number;
  tag_id: number;
  created_at: Date;
  updated_at: Date;
};

export type TaskTagModelField =
  | "task_id"
  | "tag_id"
  | "created_at"
  | "updated_at";
