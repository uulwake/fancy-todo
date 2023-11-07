export type PageInfoType = {
  size: number;
  number: number;
  total: number;
};

export type PageInfoReqQueryType = {
  page_size?: string;
  page_number?: string;
  sort_key?: string;
  sort_order?: string;
};

export type PageInfoQueryType = {
  page_size: number;
  page_number: number;
  page_offset: number;
  sort_key?: string;
  sort_order?: string;
};
