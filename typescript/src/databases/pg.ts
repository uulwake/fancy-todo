import knex, { Knex } from "knex";

export default (): Knex => {
  const instance = knex({
    client: "pg",
    connection: process.env.DATABASE_URL,
    pool: {
      min: 2,
      max: 10,
    },
  });

  // somehow afterCreate does not work https://github.com/knex/knex/issues/5352
  // workaround for checking connection
  instance.raw("SELECT 1+1 AS result").catch((err) => {
    console.error(`PostgreSQL connection error.`, err);
    process.exit(1);
  });

  return instance;
};
