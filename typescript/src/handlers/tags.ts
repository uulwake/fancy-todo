import { NextFunction, Request, Response, Router } from "express";
import { ServiceType } from "../services/types";
import { BaseHandler } from "./base";
import { authenticateJwt } from "../middlewares";
import { TagModel } from "../models/tag";
import { getUserId } from "./utils";
import { TagHandlerValidatorType } from "./validators/types";
import { createContext } from "../libs/create-context";

export class TagHandler extends BaseHandler {
  private service: ServiceType;
  private validator: TagHandlerValidatorType;

  constructor(service: ServiceType, validator: TagHandlerValidatorType) {
    super(Router());
    this.service = service;
    this.validator = validator;

    this.router.use(authenticateJwt());
    this.router.post("/", this.validator.createTag, this.createTag.bind(this));
    this.router.patch(
      "/:tag_id/tasks/:task_id",
      this.validator.addExistingTagToTask,
      this.addExistingTagToTask.bind(this)
    );
    this.router.get(
      "/search",
      this.validator.searchTag,
      this.searchTag.bind(this)
    );
    this.router.delete(
      "/:tag_id",
      this.validator.deleteTag,
      this.deleteTag.bind(this)
    );
  }

  async createTag(
    req: Request<{}, {}, Pick<TagModel, "name"> & { task_id?: number }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const tagId = await this.service.tagService.createTag(
        createContext(req),
        getUserId(req),
        req.body
      );
      res.json({ data: { tag: { id: tagId } } });
    } catch (err) {
      next(err);
    }
  }

  async addExistingTagToTask(
    req: Request<{ tag_id: string; task_id: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const tagId = Number(req.params.tag_id);
      const taskId = Number(req.params.task_id);

      await this.service.tagService.addExistingTagToTask(
        createContext(req),
        tagId,
        taskId
      );
      res.json({
        data: {
          tag: { id: tagId, task: { id: taskId } },
        },
      });
    } catch (err) {
      next(err);
    }
  }

  async searchTag(
    req: Request<{}, {}, {}, { name: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const tags = await this.service.tagService.searchTask(
        createContext(req),
        getUserId(req),
        req.query
      );
      res.json({ data: { tags } });
    } catch (err) {
      next(err);
    }
  }

  async deleteTag(
    req: Request<{ tag_id: string }>,
    res: Response,
    next: NextFunction
  ) {
    try {
      const tagId = Number(req.params.tag_id);
      await this.service.tagService.deleteTag(
        createContext(req),
        getUserId(req),
        tagId
      );
      res.json({ data: { tag: { id: tagId } } });
    } catch (err) {
      next(err);
    }
  }
}
