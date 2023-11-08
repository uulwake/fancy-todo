import { Client } from "@elastic/elasticsearch";

export default (): Client => {
  const client = new Client({
    node: process.env.ES_URL,
  });

  client.ping().catch((err) => {
    // only exit in production
    // in dev, sometimes we do not need to turn on ES
    if (process.env.NODE_ENV === "production") {
      console.error(`ElasticSearch connection error.`, err);
      process.exit(1);
    }
  });

  return client;
};
