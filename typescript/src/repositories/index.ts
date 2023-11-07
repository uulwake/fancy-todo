import { DBType } from "../databases";
import { RepositoryType } from "./types";
import { TagRepo } from "./tags";
import { TaskRepo } from "./tasks";
import { UserRepo } from "./users";

export default (db: DBType): RepositoryType => {
  const userRepo = new UserRepo(db);
  const taskRepo = new TaskRepo(db);
  const tagRepo = new TagRepo(db);

  return {
    userRepo,
    taskRepo,
    tagRepo,
  };
};
