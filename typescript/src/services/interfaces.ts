import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { UserModel } from "../models/user";
import { PageInfoQueryType } from "../types/page";

export interface IUserService {
  register(
    body: Pick<UserModel, "name" | "email" | "password">
  ): Promise<{ id: number; jwt_token: string }>;
  login(
    body: Pick<UserModel, "email" | "password">
  ): Promise<{ id: number; jwt_token: string }>;
}

export interface ITaskService {
  createTask(
    userId: number,
    data: Pick<TaskModel, "title" | "description"> & { tag_ids: number[] }
  ): Promise<number>;
  getListOfTasks(
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]>;
  getTotalOfTasks(
    userId: number,
    queryParam: { status?: string; tag_id?: number }
  ): Promise<number>;
  searchTask(
    userId: number,
    queryParam: { title: string }
  ): Promise<Pick<TaskModel, "id" | "user_id" | "title" | "status">[]>;
  getTaskDetail(userId: number, taskId: number): Promise<TaskModel | null>;
  updateTask(
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ): Promise<void>;
  deleteTask(userId: number, taskId: number): Promise<void>;
}

export interface ITagService {
  createTag(
    userId: number,
    body: Pick<TagModel, "name"> & { task_id?: number }
  ): Promise<number>;
  addExistingTagToTask(tagId: number, taskId: number): Promise<void>;
  searchTask(
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]>;
  deleteTag(userId: number, tagId: number): Promise<void>;
}
