import { Client } from "@elastic/elasticsearch";

export default () => {
  return new Client({
    node: "http://localhost:9200",
  });
};
