import { body, param, query } from "express-validator";
import { TagHandlerValidatorType } from "./types";

export const tagValidator: TagHandlerValidatorType = {
  createTag: [
    body("name")
      .notEmpty()
      .withMessage("name is empty")
      .bail()
      .isString()
      .withMessage("name is not string")
      .bail(),
    body("task_id")
      .optional()
      .isNumeric()
      .withMessage("task_id is not numeric"),
  ],
  addExistingTagToTask: [
    param("tag_id").notEmpty().withMessage("tag_id is empty"),
    param("task_id").notEmpty().withMessage("tag_id is empty"),
  ],
  searchTag: [query("name").notEmpty().withMessage("name is empty")],
  deleteTag: [param("tag_id").notEmpty().withMessage("tag_id is empty")],
};
