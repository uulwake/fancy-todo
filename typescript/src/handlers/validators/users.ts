import { body } from "express-validator";
import { UserHandlerValidatorType } from "./types";

export const userValidators: UserHandlerValidatorType = {
  register: [
    body("name")
      .notEmpty()
      .withMessage("name is empty")
      .bail()
      .isString()
      .withMessage("name is not string")
      .bail(),
    body("email")
      .notEmpty()
      .withMessage("email is empty")
      .bail()
      .isEmail()
      .withMessage("email is invalid")
      .bail(),
    body("password")
      .notEmpty()
      .withMessage("password is empty")
      .bail()
      .isString()
      .withMessage("password is not string")
      .bail(),
  ],
  login: [
    body("email")
      .notEmpty()
      .withMessage("email is empty")
      .bail()
      .isEmail()
      .withMessage("email is invalid")
      .bail(),
    body("password")
      .notEmpty()
      .withMessage("password is empty")
      .bail()
      .isString()
      .withMessage("password is not string")
      .bail(),
  ],
};
