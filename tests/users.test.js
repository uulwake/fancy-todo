const axios = require("axios");

describe("User Test", () => {
  describe("Register", () => {
    it("Should register user", async () => {
      const now = Date.now();
      const { data } = await axios({
        url: "http://localhost:3001/v1/users/register",
        method: "POST",
        data: {
          name: "foobar" + now,
          email: "foobar" + now + "@gmail.com",
          password: "12345",
        },
      });

      expect(data.data.user.id > 0).toEqual(true);
      expect(data.data.jwt_toke != "").toEqual(true);
    });

    it("Should not register user if email is invalid", async () => {
      const now = Date.now();
      try {
        await axios({
          url: "http://localhost:3001/v1/users/register",
          method: "POST",
          data: {
            name: "foobar" + now,
            email: "foobar" + now + "gmail.com",
            password: "12345",
          },
        });
      } catch (err) {
        expect(err.response.status != 200).toEqual(true);
      }
    });

    it("Should not register user if email already exists", async () => {
      const now = Date.now();
      const { data } = await axios({
        url: "http://localhost:3001/v1/users/register",
        method: "POST",
        data: {
          name: "foobar" + now,
          email: "foobar" + now + "@gmail.com",
          password: "12345",
        },
      });
      expect(data.data.user.id > 0).toEqual(true);
      expect(data.data.jwt_toke != "").toEqual(true);

      try {
        await axios({
          url: "http://localhost:3001/v1/users/register",
          method: "POST",
          data: {
            name: "foobar" + now,
            email: "foobar" + now + "@gmail.com",
            password: "12345",
          },
        });
      } catch (err) {
        expect(err.response.status != 200).toEqual(true);
      }
    });
  });

  describe("Login", () => {
    it("Should not login if password is invalid", async () => {
      try {
        const now = Date.now();
        await axios({
          url: "http://localhost:3001/v1/users/register",
          method: "POST",
          data: {
            name: "foobar" + now,
            email: "foobar" + now + "@gmail.com",
            password: "12345",
          },
        });

        await axios({
          url: "http://localhost:3001/v1/users/login",
          method: "POST",
          data: {
            email: "foobar" + now + "@gmail.com",
            password: "12345xxx",
          },
        });
      } catch (err) {
        expect(err.response.status != 200).toEqual(true);
      }
    });

    it("Should not login if email is invalid", async () => {
      try {
        const now = Date.now();
        await axios({
          url: "http://localhost:3001/v1/users/register",
          method: "POST",
          data: {
            name: "foobar" + now,
            email: "foobar" + now + "@gmail.com",
            password: "12345",
          },
        });

        await axios({
          url: "http://localhost:3001/v1/users/login",
          method: "POST",
          data: {
            email: "foobarxxx" + now + "@gmail.com",
            password: "12345",
          },
        });
      } catch (err) {
        expect(err.response.status != 200).toEqual(true);
      }
    });

    it("Should login user", async () => {
      const now = Date.now();
      const { data: dataRegister } = await axios({
        url: "http://localhost:3001/v1/users/register",
        method: "POST",
        data: {
          name: "foobar" + now,
          email: "foobar" + now + "@gmail.com",
          password: "12345",
        },
      });

      const { data: dataLogin } = await axios({
        url: "http://localhost:3001/v1/users/login",
        method: "POST",
        data: {
          email: "foobar" + now + "@gmail.com",
          password: "12345",
        },
      });

      expect(dataLogin.data.user.id).toEqual(dataRegister.data.user.id);
      expect(dataLogin.data.jwt_token != "").toEqual(true);
    });
  });
});
