import { reqValidationResult } from "../../middlewares";
import { ValidatorsType } from "./types";
import { userValidators } from "./users";
import { taskValidators } from "./tasks";
import { tagValidator } from "./tags";

const addRequestValidatorMiddleware = <T>(validator: T): T => {
  for (const key in validator) {
    (validator[key] as unknown as ValidatorsType).push(reqValidationResult());
  }

  return validator;
};

export const validators = {
  users: addRequestValidatorMiddleware(userValidators),
  tasks: addRequestValidatorMiddleware(taskValidators),
  tags: addRequestValidatorMiddleware(tagValidator),
};
