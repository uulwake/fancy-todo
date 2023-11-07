import { NextFunction, Request, Response, Router } from "express";

import { BaseHandler } from "./base";
import { authenticateJwt } from "../middlewares/authenticate";
import { ServiceType } from "../services/types";
import { TaskModel } from "../models/task";
import { PageInfoReqQueryType } from "../types/page";
import { getUserId, sanitizePageQuery } from "./utils";
import { TaskHandlerValidatorType } from "./validators/types";

export class TaskHandler extends BaseHandler {
  private service: ServiceType;
  private validator: TaskHandlerValidatorType;

  constructor(service: ServiceType, validator: TaskHandlerValidatorType) {
    super(Router());
    this.service = service;
    this.validator = validator;

    this.router.use(authenticateJwt);
    this.router.post(
      "/",
      this.validator.createTask,
      this.createTask.bind(this)
    );
    this.router.get(
      "/",
      this.validator.getListOfTasks,
      this.getListOfTasks.bind(this)
    );
    this.router.get(
      "/search",
      this.validator.searchTask,
      this.searchTask.bind(this)
    );
    this.router.get(
      "/:task_id",
      this.validator.getTaskDetail,
      this.getTaskDetail.bind(this)
    );
    this.router.patch(
      "/:task_id",
      this.validator.updateTask,
      this.updateTask.bind(this)
    );
    this.router.delete(
      "/:task_id",
      this.validator.deleteTask,
      this.deleteTask.bind(this)
    );
  }

  async createTask(
    req: Request<
      {},
      {},
      Pick<TaskModel, "title" | "description"> & { tag_ids: number[] }
    >,
    res: Response,
    next: NextFunction
  ) {
    try {
      const taskId = await this.service.taskService.createTask(
        getUserId(req),
        req.body
      );

      res.json({
        data: {
          task: {
            id: taskId,
          },
        },
      });
    } catch (err) {
      next(err);
    }
  }

  async getListOfTasks(
    req: Request<
      {},
      {},
      {},
      PageInfoReqQueryType & { status?: string; tag_id?: string }
    >,
    res: Response,
    next: NextFunction
  ) {
    try {
      const pageQuery = sanitizePageQuery({
        page_number: req.query.page_number,
        page_size: req.query.page_size,
        sort_key: req.query.sort_key,
        sort_order: req.query.sort_order,
      });

      const userId = getUserId(req);
      const whereQuery: { status?: string; tag_id?: number } = {};
      if (req.query.status) {
        whereQuery.status = req.query.status;
      }

      if (req.query.tag_id) {
        whereQuery.tag_id = Number(req.query.tag_id);
      }

      const [tasks, total] = await Promise.all([
        this.service.taskService.getListOfTasks(userId, {
          ...pageQuery,
          ...whereQuery,
        }),
        this.service.taskService.getTotalOfTasks(userId, whereQuery),
      ]);

      res.json({
        data: { tasks },
        page: {
          size: pageQuery.page_size,
          number: pageQuery.page_number,
          total,
        },
      });
    } catch (err) {
      next(err);
    }
  }

  async searchTask(
    req: Request<{}, {}, {}, { title: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const tasks = await this.service.taskService.searchTask(
        getUserId(req),
        req.query
      );
      res.json({
        data: { tasks },
      });
    } catch (err) {
      next(err);
    }
  }

  async getTaskDetail(
    req: Request<{ task_id: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const taskId = Number(req.params.task_id);
      const task = await this.service.taskService.getTaskDetail(
        getUserId(req),
        taskId
      );
      res.json({ data: { task } });
    } catch (err) {
      next(err);
    }
  }

  async updateTask(
    req: Request<
      { task_id: string },
      {},
      Partial<
        Pick<TaskModel, "title" | "description" | "status" | "order"> & {
          tag_ids: [];
        }
      >
    >,
    res: Response,
    next: NextFunction
  ) {
    try {
      const taskId = Number(req.params.task_id);
      await this.service.taskService.updateTask(
        getUserId(req),
        taskId,
        req.body
      );
      res.json({ data: { task: { id: taskId } } });
    } catch (err) {
      next(err);
    }
  }

  async deleteTask(
    req: Request<{ task_id: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const taskId = Number(req.params.task_id);
      await this.service.taskService.deleteTask(getUserId(req), taskId);
      res.json({
        data: { task: { id: taskId } },
      });
    } catch (err) {
      next(err);
    }
  }
}
