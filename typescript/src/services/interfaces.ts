import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { UserModel } from "../models/user";
import { Context } from "../types/context";
import { PageInfoQueryType } from "../types/page";

export interface IUserService {
  register(
    ctx: Context,
    body: Pick<UserModel, "name" | "email" | "password">
  ): Promise<{ id: number; jwt_token: string }>;
  login(
    ctx: Context,
    body: Pick<UserModel, "email" | "password">
  ): Promise<{ id: number; jwt_token: string }>;
}

export interface ITaskService {
  createTask(
    ctx: Context,
    userId: number,
    data: Pick<TaskModel, "title" | "description"> & { tag_ids: number[] }
  ): Promise<number>;
  getListOfTasks(
    ctx: Context,
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]>;
  getTotalOfTasks(
    ctx: Context,
    userId: number,
    queryParam: { status?: string; tag_id?: number }
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

export interface ITagService {
  createTag(
    ctx: Context,
    userId: number,
    body: Pick<TagModel, "name"> & { task_id?: number }
  ): Promise<number>;
  addExistingTagToTask(
    ctx: Context,
    tagId: number,
    taskId: number
  ): Promise<void>;
  searchTask(
    ctx: Context,
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]>;
  deleteTag(ctx: Context, userId: number, tagId: number): Promise<void>;
}
