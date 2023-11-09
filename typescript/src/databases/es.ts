import { Client } from "@elastic/elasticsearch";

export default (): Client => {
  const client = new Client({
    node: process.env.ES_URL,
  });

  client.ping().catch((err) => {
    console.error(`ElasticSearch connection error.`, err.message);

    // exit only in prod
    if (process.env.NODE_ENV === "production") {
      console.error(err);
      process.exit(1);
    }
  });

  return client;
};
