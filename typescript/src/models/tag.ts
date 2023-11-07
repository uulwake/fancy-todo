export type TagModel = {
  id: number;
  user_id: number;
  name: string;
  created_at: Date;
  updated_at: Date;
};

export type TagModelField = "id" | "name" | "created_at" | "updated_at";
