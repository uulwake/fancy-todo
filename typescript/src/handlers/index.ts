import { Router } from "express";
import { UserHandler } from "./users";
import { ServiceType } from "../services/types";
import { TaskHandler } from "./tasks";
import { TagHandler } from "./tags";
import { validators } from "./validators";

export default (service: ServiceType): Router => {
  const userHandler = new UserHandler(service, validators.users);
  const taskHandler = new TaskHandler(service, validators.tasks);
  const tagHandler = new TagHandler(service, validators.tags);

  const router = Router();
  router.use("/users", userHandler.getRouter());
  router.use("/tasks", taskHandler.getRouter());
  router.use("/tags", tagHandler.getRouter());

  return router;
};
