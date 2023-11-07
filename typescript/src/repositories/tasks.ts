import { ES_INDEX } from "../constants";
import { ITaskRepo } from "./interfaces";
import { TagModel } from "../models/tag";
import { TaskModel } from "../models/task";
import { TaskTagModel } from "../models/task_tag";
import { DBType } from "../databases";
import { PageInfoQueryType } from "../types/page";

export class TaskRepo implements ITaskRepo {
  private db: DBType;

  constructor(db: DBType) {
    this.db = db;
  }

  async getLastOrderNumber(userId: number): Promise<number> {
    const res = await this.db
      .pg<TaskModel>("tasks")
      .select("id", "order")
      .where("user_id", userId)
      .orderBy("order", "desc")
      .first();

    if (!res) {
      return 0;
    }

    return res.order;
  }

  async createTask(
    data: Omit<TaskModel, "id">,
    tagIds: number[]
  ): Promise<number> {
    const tx = await this.db.pg.transaction();

    try {
      const taskData: Omit<TaskModel, "id"> = {
        user_id: data.user_id,
        title: data.title,
        description: data.description,
        status: data.status,
        order: data.order,
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
    userId: number,
    queryParam: PageInfoQueryType & { status?: string; tag_id?: number }
  ): Promise<Omit<TaskModel, "description">[]> {
    const query = this.db
      .pg<TaskModel[]>("tasks")
      .select(
        "tasks.id",
        "tasks.user_id",
        "tasks.title",
        "tasks.description",
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

    const tasks = await query.orderBy("tasks.id", "desc");

    const mapTagByTaskId = await this.getTagsByTaskId(
      tasks.map((task) => task.id)
    );

    for (const task of tasks) {
      task.tags = mapTagByTaskId[task.id];
    }

    return tasks;
  }

  async getTotalOfTasks(
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

    const mapTagByTaskId = await this.getTagsByTaskId([taskId]);
    return { ...res, tags: mapTagByTaskId[taskId] };
  }

  async updateTask(
    userId: number,
    taskId: number,
    data: Partial<Pick<TaskModel, "title" | "description" | "status" | "order">>
  ) {
    const updatedData: Partial<TaskModel> = { ...data, updated_at: new Date() };
    await this.db
      .pg<TaskModel>("tasks")
      .update(updatedData)
      .where("user_id", userId)
      .where("id", taskId);

    const esSources = [];
    for (const key in updatedData) {
      let value;
      if (key === "updated_at") {
        value = updatedData.updated_at!.toISOString();
      } else {
        value = updatedData[key as keyof Partial<TaskModel>];
      }
      esSources.push(`ctx._source.${key} = '${value}'`);
    }

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
  }

  async deleteTask(userId: number, taskId: number): Promise<void> {
    await this.db
      .pg<TaskModel>("tasks")
      .delete()
      .where("id", taskId)
      .where("user_id", userId);

    await this.db.es.deleteByQuery({
      index: ES_INDEX.TASKS,
      query: {
        bool: {
          must: [{ match: { id: taskId } }, { match: { user_id: userId } }],
        },
      },
    });
  }
}
