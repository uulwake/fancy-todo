import express, { Request, Response } from "express";
import dotenv from "dotenv";
import cors from "cors";
import requestIp from "request-ip";

import initDB from "./databases";
import initRepo from "./repositories";
import initService from "./services";
import initHandler from "./handlers";

import { error, requestId } from "./middlewares";

dotenv.config({
  path: process.env.NODE_ENV === "production" ? ".env.production" : ".env",
});
const app = express();
app.use(express.json());
app.use(cors());
app.use(requestIp.mw());
app.use(requestId());

const db = initDB();
const repo = initRepo(db);
const service = initService(repo);
const router = initHandler(service);

app.get("/hc", (req: Request, res: Response) => {
  res.json({ status: "ok" });
});
app.use("/v1", router);
app.use(error());

export default app;
