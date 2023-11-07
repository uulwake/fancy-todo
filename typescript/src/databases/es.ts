import { Client } from "@elastic/elasticsearch";

export default () => {
  return new Client({
    node: process.env.ES_URL,
  });
};
