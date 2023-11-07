import { TASK_STATUS } from "../constants";
import { TaskModel } from "../models/task";
import { PageInfoQueryType } from "../types/page";
import { RepositoryType } from "../repositories/types";
import { ITaskService } from "./interfaces";

export class TaskService implements ITaskService {
  private repo: RepositoryType;

  constructor(repo: RepositoryType) {
    this.repo = repo;
  }

  async createTask(
    userId: number,
    data: Pick<TaskModel, "title" | "description"> & { tag_ids: number[] }
  ): Promise<number> {
    const lastOrderNum = await this.repo.taskRepo.getLastOrderNumber(userId);
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

    return this.repo.taskRepo.createTask(taskData, data.tag_ids);
  }

  async getListOfTasks(
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]> {
    return this.repo.taskRepo.getListOfTasks(userId, queryParam);
  }

  async getTotalOfTasks(
    userId: number,
    queryParam: { status?: string; tag_id?: number }
  ): Promise<number> {
    return this.repo.taskRepo.getTotalOfTasks(userId, queryParam);
  }

  async searchTask(
    userId: number,
    queryParam: { title: string }
  ): Promise<Pick<TaskModel, "id" | "user_id" | "title" | "status">[]> {
    return this.repo.taskRepo.searchTask(userId, queryParam);
  }

  async getTaskDetail(
    userId: number,
    taskId: number
  ): Promise<TaskModel | null> {
    return this.repo.taskRepo.getTaskDetail(userId, taskId);
  }

  async updateTask(
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ) {
    await this.repo.taskRepo.updateTask(userId, taskId, data);
  }

  async deleteTask(userId: number, taskId: number) {
    await this.repo.taskRepo.deleteTask(userId, taskId);
  }
}
