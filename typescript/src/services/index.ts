import { RepositoryType } from "../repositories/types";
import { ServiceType } from "./types";
import { TagService } from "./tags";
import { TaskService } from "./tasks";
import { UserService } from "./users";

export default (repo: RepositoryType): ServiceType => {
  const userService = new UserService(repo);
  const taskService = new TaskService(repo);
  const tagService = new TagService(repo);

  return {
    userService,
    taskService,
    tagService,
  };
};
