import { TASK_STATUS } from "../constants";
import { TaskModel } from "../models/task";
import { PageInfoQueryType } from "../types/page";
import { RepositoryType } from "../repositories/types";
import { ITaskService } from "./interfaces";
import { Context } from "../types/context";

export class TaskService implements ITaskService {
  private repo: RepositoryType;

  constructor(repo: RepositoryType) {
    this.repo = repo;
  }

  async createTask(
    ctx: Context,
    userId: number,
    data: Pick<TaskModel, "title" | "description"> & { tag_ids: number[] }
  ): Promise<number> {
    const lastOrderNum = await this.repo.taskRepo.getLastOrderNumber(
      ctx,
      userId
    );
    const now = new Date();
    const taskData = {
      user_id: userId,
      title: data.title,
      description: data.description,
      status: TASK_STATUS.ON_GOING,
      order: lastOrderNum + 1,
      created_at: now,
      updated_at: now,
    };

    return this.repo.taskRepo.createTask(ctx, taskData, data.tag_ids);
  }

  async getListOfTasks(
    ctx: Context,
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]> {
    return this.repo.taskRepo.getListOfTasks(ctx, userId, queryParam);
  }

  async getTotalOfTasks(
    ctx: Context,
    userId: number,
    queryParam: { status?: string; tag_id?: number }
  ): Promise<number> {
    return this.repo.taskRepo.getTotalOfTasks(ctx, userId, queryParam);
  }

  async searchTask(
    ctx: Context,
    userId: number,
    queryParam: { title: string }
  ): Promise<Pick<TaskModel, "id" | "user_id" | "title" | "status">[]> {
    return this.repo.taskRepo.searchTask(ctx, userId, queryParam);
  }

  async getTaskDetail(
    ctx: Context,
    userId: number,
    taskId: number
  ): Promise<TaskModel | null> {
    return this.repo.taskRepo.getTaskDetail(ctx, userId, taskId);
  }

  async updateTask(
    ctx: Context,
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ) {
    await this.repo.taskRepo.updateTask(ctx, userId, taskId, data);
  }

  async deleteTask(ctx: Context, userId: number, taskId: number) {
    await this.repo.taskRepo.deleteTask(ctx, userId, taskId);
  }
}
