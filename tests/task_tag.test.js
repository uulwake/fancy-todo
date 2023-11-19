const axios = require("axios");

const delay = (time) => new Promise((res) => setTimeout(res, time));
describe("Task & Tag Test", () => {
  let userId;
  let jwtToken;
  let now;
  beforeEach(async () => {
    now = Date.now();
    const { data } = await axios({
      url: "http://localhost:3001/v1/users/register",
      method: "POST",
      data: {
        name: "John Doe" + " " + now,
        email: "johndoes" + now + "@gmail.com",
        password: "secret",
      },
    });
    userId = data.data.user.id;
    jwtToken = data.data.jwt_token;
  });

  describe("Create & Get Detail", () => {
    it("Should create task with tag ID", async () => {
      let { data } = await axios({
        url: "http://localhost:3001/v1/tags",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          name: "name_" + now,
        },
      });
      const tagId = data.data.tag.id;
      expect(tagId > 0).toEqual(true);

      ({ data } = await axios({
        url: "http://localhost:3001/v1/tasks",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          title: "title_" + now,
          description: "description_" + now,
          tag_ids: [tagId],
        },
      }));

      const taskId = data.data.task.id;
      expect(taskId > 0).toEqual(true);

      ({ data } = await axios({
        url: "http://localhost:3001/v1/tasks" + "/" + taskId,
        method: "GET",
        headers: { "jwt-token": jwtToken },
      }));

      const task = data.data.task;
      expect(task.id).toEqual(taskId);
      expect(task.title).toEqual("title_" + now);
      expect(task.description).toEqual("description_" + now);
      expect(task.status).toEqual("on_going");
      expect(task.order).toEqual(null);
      expect(task.tags[0].id).toEqual(tagId);
      expect(task.tags[0].name).toEqual("name_" + now);
    });

    it("Should create tag with task ID", async () => {
      let { data } = await axios({
        url: "http://localhost:3001/v1/tasks",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          title: "title_" + now,
          description: "description_" + now,
        },
      });

      const taskId = data.data.task.id;
      expect(taskId > 0).toEqual(true);

      ({ data } = await axios({
        url: "http://localhost:3001/v1/tags",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          name: "name_" + now,
          task_id: taskId,
        },
      }));

      const tagId = data.data.tag.id;
      expect(tagId > 0).toEqual(true);

      ({ data } = await axios({
        url: "http://localhost:3001/v1/tasks" + "/" + taskId,
        method: "GET",
        headers: { "jwt-token": jwtToken },
      }));

      const task = data.data.task;
      expect(task.id).toEqual(taskId);
      expect(task.title).toEqual("title_" + now);
      expect(task.description).toEqual("description_" + now);
      expect(task.status).toEqual("on_going");
      expect(task.order).toEqual(null);
      expect(task.tags[0].id).toEqual(tagId);
      expect(task.tags[0].name).toEqual("name_" + now);
    });
  });

  describe("Lists, Search, Update, Delete", () => {
    let taskId;
    let tagId;
    beforeEach(async () => {
      let { data } = await axios({
        url: "http://localhost:3001/v1/tasks",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          title: "title_" + now,
          description: "description_" + now,
        },
      });
      taskId = data.data.task.id;

      ({ data } = await axios({
        url: "http://localhost:3001/v1/tags",
        method: "POST",
        headers: { "jwt-token": jwtToken },
        data: {
          name: "name_" + now,
          task_id: taskId,
        },
      }));
      tagId = data.data.tag.id;
    });

    describe("List", () => {
      it("Should get lists", async () => {
        let { data } = await axios({
          url: "http://localhost:3001/v1/tasks",
          method: "GET",
          headers: { "jwt-token": jwtToken },
        });

        expect(data.page.size).toEqual(10);
        expect(data.page.total).toEqual(1);
        expect(data.data.tasks.length).toEqual(1);
        expect(data.data.tasks[0].id).toEqual(taskId);
        expect(data.data.tasks[0].tags[0].id).toEqual(tagId);
      });
    });

    describe.only("Search", () => {
      let taskIds;
      let tagIds;
      beforeEach(async () => {
        const promiseTasks = ["abc", "bcd", "cde"].map(async (title) =>
          axios({
            url: "http://localhost:3001/v1/tasks",
            method: "POST",
            headers: { "jwt-token": jwtToken },
            data: {
              title,
              description: "description",
            },
          })
        );

        const promiseTags = ["vwx", "wxy", "xyz"].map(async (name) =>
          axios({
            url: "http://localhost:3001/v1/tags",
            method: "POST",
            headers: { "jwt-token": jwtToken },
            data: {
              name,
            },
          })
        );

        const res = await Promise.all([...promiseTasks, ...promiseTags]);

        taskIds = res.slice(0, 3).map((el) => el.data.data.task.id);
        tagIds = res.slice(3).map((el) => el.data.data.tag.id);
      });

      afterEach(async () => {
        const promiseTasks = taskIds.map(async (id) =>
          axios({
            url: "http://localhost:3001/v1/tasks/" + id,
            method: "DELETE",
            headers: { "jwt-token": jwtToken },
          })
        );

        const promiseTags = tagIds.map(async (id) =>
          axios({
            url: "http://localhost:3001/v1/tags/" + id,
            method: "DELETE",
            headers: { "jwt-token": jwtToken },
          })
        );

        await Promise.all([...promiseTasks, ...promiseTags]);
      });

      it.only("Should search tasks and tags", async () => {
        await delay(2e3); // wait for ES

        let { data } = await axios({
          url: "http://localhost:3001/v1/tasks/search",
          params: { title: "Bc" },
          method: "GET",
          headers: { "jwt-token": jwtToken },
        });
        expect(taskIds.includes(data.data.tasks[0].id));
        expect(taskIds.includes(data.data.tasks[1].id));
        expect(["abc", "bcd"].includes(data.data.tasks[0].title));
        expect(["abc", "bcd"].includes(data.data.tasks[1].title));

        ({ data } = await axios({
          url: "http://localhost:3001/v1/tags/search",
          params: { name: "xY" },
          method: "GET",
          headers: { "jwt-token": jwtToken },
        }));
        expect(tagIds.includes(data.data.tags[0].id));
        expect(tagIds.includes(data.data.tags[1].id));
        expect(["wxy", "xyz"].includes(data.data.tags[0].name));
        expect(["wxy", "xyz"].includes(data.data.tags[1].name));
      });
    });
  });
});
