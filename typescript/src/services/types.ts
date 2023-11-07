import { ITagService, ITaskService, IUserService } from "./interfaces";

export type ServiceType = {
  userService: IUserService;
  taskService: ITaskService;
  tagService: ITagService;
};
