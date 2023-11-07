import { ITagRepo, ITaskRepo, IUserRepo } from "./interfaces";

export type RepositoryType = {
  userRepo: IUserRepo;
  taskRepo: ITaskRepo;
  tagRepo: ITagRepo;
};
