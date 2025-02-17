<script>
import { EventBus } from "../EventBus";
export default {
  data() {
    return {
      username: "",
    };
  },
  methods: {
    async handleLogin() {
      try {
        const response = await this.$axios.post("/session", {
          username: this.username,
        });
        if (response.status == 200) {
          alert("Welcome Back " + this.username);
          localStorage.setItem("username", this.username);
          localStorage.setItem("authToken", response.data["resourceId"]);
          localStorage.setItem("isLoggedIn", "true");
          EventBus.isLoggedIn = true;
          this.$router.replace("/");
        } else if (response.status == 201) {
          alert("Welcome " + this.username);
          localStorage.setItem("username", this.username);
          localStorage.setItem("authToken", response.data["resourceId"]);
          localStorage.setItem("isLoggedIn", "true");
          EventBus.isLoggedIn = true;
          this.$router.replace("/");
        }
      } catch (error) {
        this.error = "Failed to load data: " + error.message;
      }
    },
  },
};
</script>

<template>
  <div class="login-body">
    <div class="login-form">
      <h3 class="login-title">Welcome Back!</h3>
      <div class="form-group">
        <input
          type="text"
          id="username"
          placeholder="Username"
          v-model="username"
          required
        />
      </div>
      <button class="login-button" @click="handleLogin()">
        <b>Login</b>
      </button>
    </div>
  </div>
</template>

<style>
.login-body {
  height: 80vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f6f6f6;
}

.login-form {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  width: 250px;
  padding: 25px;
  background-color: #fff;
  border-radius: 25px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.login-title {
  margin: 0;
  color: #ff6f61;
  font-family: "Arial", sans-serif;
  font-size: 22px;
  font-weight: bold;
}

.form-group {
  margin: 20px 0;
}

input {
  width: 100%;
  padding: 12px;
  margin: 8px 0;
  border: 2px solid #ccc;
  border-radius: 15px;
  font-size: 12px;
  outline: none;
}

.login-button {
  padding: 12px 20px;
  width: 100%;
  background-color: #ff6f61;
  color: white;
  border: none;
  border-radius: 20px;
  font-weight: bold;
  font-size: 14px;
}
</style>
