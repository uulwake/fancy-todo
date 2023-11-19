import { body, param, query } from "express-validator";
import { TaskHandlerValidatorType } from "./types";

export const taskValidators: TaskHandlerValidatorType = {
  createTask: [
    body("title")
      .notEmpty()
      .withMessage("title is empty")
      .bail()
      .isString()
      .withMessage("title is not string")
      .bail(),
    body("description")
      .notEmpty()
      .withMessage("description is empty")
      .isString()
      .withMessage("description is not string")
      .bail(),
    body("tag_ids")
      .optional()
      .notEmpty()
      .withMessage("tag_ids is empty")
      .bail()
      .isArray()
      .withMessage("tag_ids is not array")
      .bail(),
  ],
  getListOfTasks: [
    query("page_size").optional(),
    query("page_number").optional(),
    query("sort_key").optional(),
    query("sort_order").optional(),
    query("status").optional(),
    query("tag_id").optional(),
  ],
  searchTask: [query("title").notEmpty().withMessage("title is empty")],
  getTaskDetail: [param("task_id").notEmpty().withMessage("task_id is empty")],
  updateTask: [
    param("task_id").notEmpty().withMessage("task_id is empty"),
    body("title").optional().isString().withMessage("title is not string"),
    body("description")
      .optional()
      .isString()
      .withMessage("description is not string"),
    body("status").optional().isString().withMessage("status is not string"),
    body("order").optional().isNumeric().withMessage("order is not numeric"),
    body("tag_ids").optional().isArray().withMessage("tag_ids is not an array"),
  ],
  deleteTask: [param("task_id").notEmpty().withMessage("task_id is empty")],
};
