import { Client } from "@elastic/elasticsearch";
import { Knex } from "knex";

import pgInit from "./pg";
import esInit from "./es";

export type DBType = {
  pg: Knex;
  es: Client;
};


export default (): DBType => {
  const pg = pgInit();
  const es = esInit();

  return {
    pg,
    es,
  };
};
