import jwt from "jsonwebtoken";
import bcrypt from "bcrypt";
import { ITagRepo, ITaskRepo, IUserRepo } from "../repositories/interfaces";
import { UserService } from "./users";
import { mock } from "jest-mock-extended";

const mockUserRepo = mock<IUserRepo>();
const mockTaskRepo = mock<ITaskRepo>();
const mockTagRepo = mock<ITagRepo>();
const userService = new UserService({
  userRepo: mockUserRepo,
  taskRepo: mockTaskRepo,
  tagRepo: mockTagRepo,
});

jest.mock("jsonwebtoken");
jest.mock("bcrypt");

describe("User Service", () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe("Register", () => {
    it("Should register user", async () => {
      (jwt.sign as any).mockReturnValue("token");
      (bcrypt.hash as any).mockResolvedValue("password");
      jest.useFakeTimers().setSystemTime(new Date("2023-01-01"));
      mockUserRepo.createUser.mockResolvedValue(1);

      const result = await userService.register(
        { ipAddr: "ip", requestId: "reqId" },
        { name: "name", email: "email", password: "secret" }
      );

      expect(result).toEqual({ id: 1, jwt_token: "token" });

      expect(jwt.sign).toHaveBeenCalledTimes(1);
      expect(bcrypt.hash).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.createUser).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.createUser).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "reqId",
        },
        {
          name: "name",
          email: "email",
          password: "password",
          created_at: new Date("2023-01-01"),
          updated_at: new Date("2023-01-01"),
        }
      );
    });
  });

  describe("Login", () => {
    it("Should not login user if user or user ID or email or password do not exist", async () => {
      mockUserRepo.getUserDetail.mockResolvedValue({});
      await expect(
        userService.login(
          { ipAddr: "ip", requestId: "reqId" },
          { email: "email", password: "secret" }
        )
      ).rejects.toThrow("invalid email/password");

      expect(mockUserRepo.getUserDetail).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.getUserDetail).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "reqId",
        },
        { email: "email", cols: ["id", "email", "password"] }
      );
    });

    it("Should not login user if compare sync is error", async () => {
      mockUserRepo.getUserDetail.mockResolvedValue({
        id: 1,
        email: "email",
        password: "password",
      });
      (bcrypt.compare as any).mockResolvedValue(false);

      await expect(
        userService.login(
          { ipAddr: "ip", requestId: "reqId" },
          { email: "email", password: "secret" }
        )
      ).rejects.toThrow("invalid email/password");

      expect(mockUserRepo.getUserDetail).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.getUserDetail).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "reqId",
        },
        { email: "email", cols: ["id", "email", "password"] }
      );

      expect(bcrypt.compare).toHaveBeenCalledTimes(1);
    });

    it("Should login user", async () => {
      mockUserRepo.getUserDetail.mockResolvedValue({
        id: 1,
        email: "email",
        password: "password",
      });
      (bcrypt.compare as any).mockResolvedValue(true);
      (jwt.sign as any).mockReturnValue("token");

      const result = await userService.login(
        { ipAddr: "ip", requestId: "reqId" },
        { email: "email", password: "secret" }
      );

      expect(result).toEqual({ id: 1, jwt_token: "token" });
      expect(jwt.sign).toHaveBeenCalledTimes(1);
      expect(bcrypt.compare).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.getUserDetail).toHaveBeenCalledTimes(1);
      expect(mockUserRepo.getUserDetail).toHaveBeenCalledWith(
        {
          ipAddr: "ip",
          requestId: "reqId",
        },
        { email: "email", cols: ["id", "email", "password"] }
      );
    });
  });
});
