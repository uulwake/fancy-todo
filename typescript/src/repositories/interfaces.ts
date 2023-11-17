import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { TaskTagModel } from "../models/task_tag";
import { UserModel, UserModelField } from "../models/user";
import { Context } from "../types/context";
import { PageInfoQueryType } from "../types/page";

export interface IUserRepo {
  createUser(ctx: Context, data: Omit<UserModel, "id">): Promise<number>;
  getUserDetail(
    ctx: Context,
    opt: {
      id?: number;
      email?: string;
      cols?: UserModelField[];
    }
  ): Promise<Partial<UserModel>>;
}

export interface ITaskRepo {
  createTask(
    ctx: Context,
    data: Omit<TaskModel, "id" | "order">,
    tagIds: number[]
  ): Promise<number>;
  getListOfTasks(
    ctx: Context,
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]>;
  getTotalOfTasks(
    ctx: Context,
    userId: number,
    queryParam: { status?: string }
  ): Promise<number>;
  searchTask(
    ctx: Context,
    userId: number,
    queryParam: { title: string }
  ): Promise<Pick<TaskModel, "id" | "user_id" | "title" | "status">[]>;
  getTaskDetail(
    ctx: Context,
    userId: number,
    taskId: number
  ): Promise<TaskModel | null>;
  updateTask(
    ctx: Context,
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ): Promise<void>;
  deleteTask(ctx: Context, userId: number, taskId: number): Promise<void>;
}

export interface ITagRepo {
  createTag(
    ctx: Context,
    data: Omit<TagModel, "id">,
    taskId?: number
  ): Promise<number>;
  addExistingTagToTask(ctx: Context, data: TaskTagModel): Promise<void>;
  searchTag(
    ctx: Context,
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]>;
  deleteTag(ctx: Context, userId: number, tagId: number): Promise<void>;
}
