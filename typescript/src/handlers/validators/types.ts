import { RequestHandler } from "express";
import { ValidationChain } from "express-validator";

export type HandlerValidatorType = {
  users: UserHandlerValidatorType;
  tasks: TaskHandlerValidatorType;
  tags: TagHandlerValidatorType;
};

export type ValidatorsType = (ValidationChain | RequestHandler)[];

export type UserHandlerValidatorType = {
  register: ValidatorsType;
  login: ValidatorsType;
};

export type TaskHandlerValidatorType = {
  createTask: ValidatorsType;
  getListOfTasks: ValidatorsType;
  searchTask: ValidatorsType;
  getTaskDetail: ValidatorsType;
  updateTask: ValidatorsType;
  deleteTask: ValidatorsType;
};

export type TagHandlerValidatorType = {
  createTag: ValidatorsType;
  addExistingTagToTask: ValidatorsType;
  searchTag: ValidatorsType;
  deleteTag: ValidatorsType;
};
