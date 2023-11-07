import { Router } from "express";

export class BaseHandler {
  protected router: Router;
  constructor(router: Router) {
    this.router = router;
  }

  getRouter(): Router {
    return this.router;
  }
}
