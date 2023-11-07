import knex, { Knex } from "knex";

export default (): Knex => {
  return knex({
    client: "pg",
    connection: process.env.DATABASE_URL,
  });
};
