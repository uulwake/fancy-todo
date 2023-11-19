import { TagModel } from "../models/tag";
import { RepositoryType } from "../repositories/types";
import { Context } from "../types/context";
import { ITagService } from "./interfaces";

export class TagService implements ITagService {
  private repo: RepositoryType;

  constructor(repo: RepositoryType) {
    this.repo = repo;
  }

  async createTag(
    ctx: Context,
    userId: number,
    body: Pick<TagModel, "name"> & { task_id?: number }
  ): Promise<number> {
    const now = new Date();

    const data: Omit<TagModel, "id"> = {
      user_id: userId,
      name: body.name,
      created_at: now,
      updated_at: now,
    };

    return this.repo.tagRepo.createTag(ctx, data, body.task_id);
  }

  async addExistingTagToTask(
    ctx: Context,
    userId: number,
    tagId: number,
    taskId: number
  ) {
    const now = new Date();
    await this.repo.tagRepo.addExistingTagToTask(ctx, userId, {
      tag_id: tagId,
      task_id: taskId,
      created_at: now,
      updated_at: now,
    });
  }

  async searchTask(
    ctx: Context,
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]> {
    return this.repo.tagRepo.searchTag(ctx, userId, queryParam);
  }

  async deleteTag(ctx: Context, userId: number, tagId: number) {
    await this.repo.tagRepo.deleteTag(ctx, userId, tagId);
  }
}
