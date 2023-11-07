import { Request } from "express";
import { PAGE_DEFAULT } from "../constants";
import { PageInfoQueryType, PageInfoReqQueryType } from "../types/page";
import { CustomError } from "../libs/custom-error";

export const getUserId = (req: Request): number => {
  if (!req.user_id)
    throw new CustomError({
      message: "Missing user ID in requests",
      status: 401,
    });

  return req.user_id;
};

export const sanitizePageQuery = (
  page: PageInfoReqQueryType
): PageInfoQueryType => {
  let cleanPageNumber: number, cleanPageSize: number;

  if (!page.page_number) {
    cleanPageNumber = PAGE_DEFAULT.NUMBER;
  } else {
    const pageNumber = Number(page.page_number);
    if (pageNumber < 0 || pageNumber > PAGE_DEFAULT.MAX_NUMBER) {
      cleanPageNumber = PAGE_DEFAULT.MAX_NUMBER;
    } else {
      cleanPageNumber = pageNumber;
    }
  }

  if (!page.page_size) {
    cleanPageSize = PAGE_DEFAULT.SIZE;
  } else {
    const pageSize = Number(page.page_size);
    if (pageSize < 0 || pageSize > PAGE_DEFAULT.MAX_SIZE) {
      cleanPageSize = PAGE_DEFAULT.MAX_SIZE;
    } else {
      cleanPageSize = pageSize;
    }
  }

  return {
    page_number: cleanPageNumber,
    page_size: cleanPageSize,
    page_offset: cleanPageNumber * cleanPageSize,
    sort_key: page.sort_key,
    sort_order: page.sort_order,
  };
};
