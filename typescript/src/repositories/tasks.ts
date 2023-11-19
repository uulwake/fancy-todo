import { ES_INDEX } from "../constants";
import { ITaskRepo } from "./interfaces";
import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { TaskTagModel } from "../models/task_tag";
import { DBType } from "../databases";
import { PageInfoQueryType } from "../types/page";
import { Context } from "../types/context";

export class TaskRepo implements ITaskRepo {
  private db: DBType;

  constructor(db: DBType) {
    this.db = db;
  }

  async createTask(
    ctx: Context,
    data: Omit<TaskModel, "id" | "order">,
    tagIds: number[] = []
  ): Promise<number> {
    const tx = await this.db.pg.transaction();

    try {
      const taskData: Omit<TaskModel, "id" | "order"> = {
        user_id: data.user_id,
        title: data.title,
        description: data.description,
        status: data.status,
        created_at: data.created_at,
        updated_at: data.updated_at,
      };

      const res = await tx<TaskModel>("tasks").insert(taskData).returning("id");
      const taskId = res[0].id;

      const taskTagData: TaskTagModel[] = [];
      for (const tagId of tagIds) {
        taskTagData.push({
          tag_id: tagId,
          task_id: taskId,
          created_at: taskData.created_at,
          updated_at: taskData.updated_at,
        });
      }

      if (taskTagData.length) {
        await tx<TaskTagModel>("tasks_tags").insert(taskTagData);
      }

      await this.db.es.index({
        index: ES_INDEX.TASKS,
        document: {
          id: taskId,
          ...taskData,
          created_at: taskData.created_at.toISOString(),
          updated_at: taskData.updated_at.toISOString(),
        },
      });

      await tx.commit();
      return taskId;
    } catch (err) {
      await tx.rollback();
      throw err;
    }
  }

  async getTagsByTaskId(
    ctx: Context,
    taskIds: number[]
  ): Promise<{ [key: number]: TagModel[] }> {
    const tags = await this.db
      .pg<TagModel>("tags")
      .select("tags.id", "tags.name", "tasks_tags.task_id")
      .leftJoin("tasks_tags", "tasks_tags.tag_id", "tags.id")
      .whereIn("tasks_tags.task_id", taskIds)
      .orderBy("tags.id", "asc");

    const mapTagByTaskId: { [key: number]: TagModel[] } = {};
    for (const tag of tags) {
      const { task_id: taskId } = tag;
      if (taskId in mapTagByTaskId) {
        mapTagByTaskId[taskId].push(tag);
      } else {
        mapTagByTaskId[taskId] = [tag];
      }
    }

    return mapTagByTaskId;
  }

  async getListOfTasks(
    ctx: Context,
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]> {
    const query = this.db
      .pg<TaskModel[]>("tasks")
      .select(
        "tasks.id",
        "tasks.user_id",
        "tasks.title",
        "tasks.status",
        "tasks.order",
        "tasks.created_at",
        "tasks.updated_at"
      )
      .where("tasks.user_id", userId)
      .limit(queryParam.page_size)
      .offset(queryParam.page_offset);

    if (queryParam.status) {
      query.where("tasks.status", queryParam.status);
    }

    if (queryParam.tag_id) {
      query
        .leftJoin("tasks_tags", "tasks_tags.task_id", "tasks.id")
        .where("tasks_tags.tag_id", queryParam.tag_id);
    }

    if (queryParam.sort_key && queryParam.sort_order) {
      query.orderBy(queryParam.sort_key, queryParam.sort_order);
    }

    query.orderByRaw(`tasks."order" ASC NULLS LAST`);
    query.orderBy("tasks.id", "ASC");
    const tasks = await query;

    const mapTagByTaskId = await this.getTagsByTaskId(
      ctx,
      tasks.map((task) => task.id)
    );

    for (const task of tasks) {
      task.tags = mapTagByTaskId[task.id];
    }

    return tasks;
  }

  async getTotalOfTasks(
    ctx: Context,
    userId: number,
    queryParam: { status?: string; tag_id?: number }
  ): Promise<number> {
    const query = this.db
      .pg("tasks")
      .count("id as total")
      .where("user_id", userId)
      .first();

    if (queryParam.status) {
      query.where("status", queryParam.status);
    }

    if (queryParam.tag_id) {
      query
        .leftJoin("tasks_tags", "tasks_tags.task_id", "tasks.id")
        .where("tasks_tags.tag_id", queryParam.tag_id);
    }

    const { total } = (await query) as any as { total: string };
    return Number(total);
  }

  async searchTask(
    ctx: Context,
    userId: number,
    queryParam: { title: string }
  ): Promise<Pick<TaskModel, "id" | "user_id" | "title" | "status">[]> {
    const res = await this.db.es.search<TaskModel>({
      index: ES_INDEX.TASKS,
      query: {
        bool: {
          must: {
            match: {
              user_id: userId,
            },
          },
          filter: {
            wildcard: {
              title: {
                value: "*" + queryParam.title + "*",
                case_insensitive: true,
              },
            },
          },
        },
      },
    });

    const tasks: Pick<TaskModel, "id" | "user_id" | "title" | "status">[] = [];

    for (const hit of res.hits.hits) {
      if (!hit._source) continue;

      tasks.push({
        id: hit._source.id,
        user_id: hit._source.user_id,
        title: hit._source.title,
        status: hit._source.status,
      });
    }

    return tasks;
  }

  async getTaskDetail(
    ctx: Context,
    userId: number,
    taskId: number
  ): Promise<(TaskModel & { tags: TagModel[] }) | null> {
    const res = await this.db
      .pg<TaskModel>("tasks")
      .select("*")
      .where("user_id", userId)
      .where("id", taskId)
      .first();

    if (!res) return null;

    const mapTagByTaskId = await this.getTagsByTaskId(ctx, [taskId]);
    return { ...res, tags: mapTagByTaskId[taskId] };
  }

  async updateTask(
    ctx: Context,
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ) {
    const tx = await this.db.pg.transaction();

    try {
      const esSources = [];
      const updatedData: Partial<TaskModel> = {};

      updatedData.updated_at = new Date();
      esSources.push(
        `ctx._source.updated_at = '${updatedData.updated_at.toISOString()}'`
      );

      if ("title" in data) {
        updatedData.title = data.title;
        esSources.push(`ctx._source.title = '${data.title}'`);
      }

      if ("description" in data) {
        updatedData.description = data.description;
        esSources.push(`ctx._source.description = '${data.description}'`);
      }

      if ("status" in data) {
        updatedData.status = data.status;
        esSources.push(`ctx._source.status = '${data.status}'`);
      }

      if ("order" in data) {
        updatedData.order = data.order;
        esSources.push(`ctx._source.order = ${data.order}`);
      }

      await tx("tasks")
        .update(updatedData)
        .where("user_id", userId)
        .where("id", taskId);

      await this.db.es.updateByQuery({
        index: ES_INDEX.TASKS,
        query: {
          bool: {
            must: [{ match: { id: taskId } }, { match: { user_id: userId } }],
          },
        },
        script: {
          lang: "painless",
          source: esSources.join(";"),
        },
      });

      await tx.commit();
    } catch (err) {
      await tx.rollback();
      throw err;
    }
  }

  async deleteTask(
    ctx: Context,
    userId: number,
    taskId: number
  ): Promise<void> {
    const tx = await this.db.pg.transaction();

    try {
      await tx("tasks").delete().where("id", taskId).where("user_id", userId);

      await this.db.es.deleteByQuery({
        index: ES_INDEX.TASKS,
        query: {
          bool: {
            must: [{ match: { id: taskId } }, { match: { user_id: userId } }],
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
