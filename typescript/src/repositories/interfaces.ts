import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { TaskTagModel } from "../models/task_tag";
import { UserModel, UserModelField } from "../models/user";
import { PageInfoQueryType } from "../types/page";

export interface IUserRepo {
  createUser(data: Omit<UserModel, "id">): Promise<number>;
  getUserDetail(opt: {
    id?: number;
    email?: string;
    cols?: UserModelField[];
  }): Promise<Partial<UserModel>>;
}

export interface ITaskRepo {
  getLastOrderNumber(userId: number): Promise<number>;
  createTask(data: Omit<TaskModel, "id">, tagIds: number[]): Promise<number>;
  getListOfTasks(
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]>;
  getTotalOfTasks(
    userId: number,
    queryParam: { status?: string }
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

export interface ITagRepo {
  createTag(data: Omit<TagModel, "id">, taskId?: number): Promise<number>;
  addExistingTagToTask(data: TaskTagModel): Promise<void>;
  searchTag(
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]>;
  deleteTag(userId: number, tagId: number): Promise<void>;
}
