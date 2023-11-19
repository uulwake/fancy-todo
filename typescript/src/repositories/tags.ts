import { ES_INDEX } from "../constants";
import { ITagRepo } from "./interfaces";
import { TagModel } from "../models/tag";
import { TaskTagModel } from "../models/task_tag";
import { DBType } from "../databases";
import { Context } from "../types/context";

export class TagRepo implements ITagRepo {
  private db: DBType;

  constructor(db: DBType) {
    this.db = db;
  }

  async createTag(
    ctx: Context,
    data: Omit<TagModel, "id">,
    taskId?: number
  ): Promise<number> {
    const tx = await this.db.pg.transaction();
    try {
      const res = await tx<TagModel>("tags").insert(data).returning("id");
      const tagId = res[0].id;

      if (taskId) {
        await tx<TaskTagModel>("tasks_tags").insert({
          task_id: taskId,
          tag_id: tagId,
          created_at: data.created_at,
          updated_at: data.updated_at,
        });
      }

      // insert to es
      await this.db.es.index({
        index: ES_INDEX.TAGS,
        document: {
          id: tagId,
          ...data,
        },
      });

      await tx.commit();
      return tagId;
    } catch (err) {
      await tx.rollback();
      throw err;
    }
  }

  async addExistingTagToTask(
    ctx: Context,
    userId: number,
    data: TaskTagModel
  ): Promise<void> {
    let res = (await this.db
      .pg("tasks")
      .count("id as total")
      .where("id", data.task_id)
      .where("user_id", userId)
      .first()) as any as { total: number };

    if (res.total != 1) {
      throw new Error(
        `user ID: ${userId} does not have task ID: ${data.task_id}`
      );
    }

    res = (await this.db
      .pg("tags")
      .count("id as total")
      .where("id", data.tag_id)
      .where("user_id", userId)
      .first()) as any as { total: number };

    if (res.total != 1) {
      throw new Error(
        `user ID: ${userId} does not have tag ID: ${data.tag_id}`
      );
    }

    await this.db.pg<TaskTagModel>("tasks_tags").insert(data);
  }

  async searchTag(
    ctx: Context,
    userId: number,
    queryParam: { name: string }
  ): Promise<Pick<TagModel, "id" | "name">[]> {
    const res = await this.db.es.search<TagModel>({
      index: ES_INDEX.TAGS,
      query: {
        bool: {
          must: {
            match: {
              user_id: userId,
            },
          },
          filter: {
            wildcard: {
              name: {
                value: "*" + queryParam.name + "*",
                case_insensitive: true,
              },
            },
          },
        },
      },
    });

    const tags: Pick<TagModel, "id" | "name">[] = [];
    for (const hit of res.hits.hits) {
      if (!hit._source) continue;

      tags.push({
        id: hit._source?.id,
        name: hit._source?.name,
      });
    }

    return tags;
  }

  async deleteTag(ctx: Context, userId: number, tagId: number): Promise<void> {
    const tx = await this.db.pg.transaction();
    try {
      await tx("tags").delete().where("user_id", userId).where("id", tagId);

      await this.db.es.deleteByQuery({
        index: ES_INDEX.TAGS,
        query: {
          bool: {
            must: [{ match: { id: tagId } }, { match: { user_id: userId } }],
          },
        },
      });
      await tx.commit();
    } catch (err) {
      await tx.rollback();
      throw err;
    }
  }
}
