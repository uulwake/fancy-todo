import { getMockReq, getMockRes } from "@jest-mock/express";

import {
  ITagService,
  ITaskService,
  IUserService,
} from "../services/interfaces";
import { UserHandler } from "./users";
import { validators } from "./validators";
import { mock } from "jest-mock-extended";

const mockUserService = mock<IUserService>();
const mockTaskService = mock<ITaskService>();
const mockTagService = mock<ITagService>();
const userHandler = new UserHandler(
  {
    userService: mockUserService,
    tagService: mockTagService,
    taskService: mockTaskService,
  },
  validators.users
);

describe("User Handler", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe("Register", () => {
    it("Should handle error", async () => {
      mockUserService.register.mockRejectedValue("server error");

      const req = getMockReq({
        body: {
          name: "name",
          email: "mail@mail.com",
          password: "secret",
        },
        headers: {
          request_id: "req-1",
        },
        clientIp: "ip",
        user_id: null,
      });

      const { res, next } = getMockRes();

      await userHandler.register(req, res, next);

      expect(mockUserService.register).toHaveBeenCalledTimes(1);
      expect(mockUserService.register).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "req-1",
          userId: null,
        },
        req.body
      );

      expect(res.json).toHaveBeenCalledTimes(0);

      expect(next).toHaveBeenCalledTimes(1);
      expect(next).toHaveBeenCalledWith("server error");
    });

    it("Should register user", async () => {
      mockUserService.register.mockResolvedValue({ id: 1, jwt_token: "token" });

      const req = getMockReq({
        body: {
          name: "name",
          email: "email@google.com",
          password: "secret",
        },
        headers: {
          request_id: "req-1",
        },
        clientIp: "ip",
        user_id: null,
      });

      const { res, next } = getMockRes();

      await userHandler.register(req, res, next);

      expect(mockUserService.register).toHaveBeenCalledTimes(1);
      expect(mockUserService.register).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "req-1",
          userId: null,
        },
        req.body
      );

      expect(res.json).toHaveBeenCalledTimes(1);
      expect(res.json).toHaveBeenCalledWith({
        data: { user: { id: 1 }, jwt_token: "token" },
      });

      expect(next).toHaveBeenCalledTimes(0);
    });
  });

  describe("Login", () => {
    it("Should handle error", async () => {
      mockUserService.login.mockRejectedValue("server error");

      const req = getMockReq({
        body: {
          email: "mail@mail.com",
          password: "secret",
        },
        headers: {
          request_id: "req-1",
        },
        clientIp: "ip",
        user_id: null,
      });

      const { res, next } = getMockRes();

      await userHandler.login(req, res, next);

      expect(mockUserService.login).toHaveBeenCalledTimes(1);
      expect(mockUserService.login).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "req-1",
          userId: null,
        },
        req.body
      );

      expect(res.json).toHaveBeenCalledTimes(0);

      expect(next).toHaveBeenCalledTimes(1);
      expect(next).toHaveBeenCalledWith("server error");
    });

    it("Should login user", async () => {
      mockUserService.login.mockResolvedValue({
        id: 1,
        jwt_token: "token",
      });

      const req = getMockReq({
        body: {
          email: "email@google.com",
          password: "secret",
        },
        headers: {
          request_id: "req-1",
        },
        clientIp: "ip",
        user_id: null,
      });

      const { res, next } = getMockRes();

      await userHandler.login(req, res, next);

      expect(mockUserService.login).toHaveBeenCalledTimes(1);
      expect(mockUserService.login).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "req-1",
          userId: null,
        },
        req.body
      );

      expect(res.json).toHaveBeenCalledTimes(1);
      expect(res.json).toHaveBeenCalledWith({
        data: { user: { id: 1 }, jwt_token: "token" },
      });

      expect(next).toHaveBeenCalledTimes(0);
    });
  });
});
