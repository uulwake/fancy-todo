import express, { Request, Response } from "express";
import dotenv from "dotenv";
import cors from "cors";

import initDB from "./databases";
import initRepo from "./repositories";
import initService from "./services";
import initHandler from "./handlers";

import error from "./middlewares/error";

dotenv.config({
  path: process.env.NODE_ENV === "production" ? ".env.production" : ".env",
});
const app = express();
app.use(express.json());
app.use(cors());

const db = initDB();
const repo = initRepo(db);
const service = initService(repo);
const router = initHandler(service);

app.get("/healthcheck", (req: Request, res: Response) => {
  res.json({ status: "ok" });
});
app.use("/v1", router);
app.use(error);

export default app;
