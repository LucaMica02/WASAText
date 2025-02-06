<script setup>
import { RouterLink, RouterView } from "vue-router";
import { EventBus } from "./EventBus";
</script>
<script>
export default {
  data() {
    return {
      isLoggedIn: localStorage.getItem("isLoggedIn") === "true",
    };
  },
  methods: {
    logout() {
      localStorage.setItem("username", "");
      localStorage.setItem("authToken", "");
      localStorage.setItem("isLoggedIn", "false");
      EventBus.isLoggedIn = false;
      this.$router.replace("/login");
    },
    getUsers() {
      this.$router.replace("/users");
    },
    getProfile() {
      this.$router.replace("/profile");
    },
    createNewGroup() {
      this.$router.replace("/newGroup");
    },
    login() {
      this.$router.replace("/login");
    },
    goHome() {
      this.$router.replace("/");
    },
  },
  watch: {
    isLoggedIn(newVal) {
      localStorage.setItem("isLoggedIn", newVal.toString());
    },
  },
};
</script>

<template>
  <header class="header">
    <span class="clickable-container" @click="goHome()">
      <img
        class="image"
        src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRc_8Zx8cH_HmDz1lfEjqGOsnnDlfhYzDz1UA&s"
        alt="image"
      /><b>WASAText</b></span
    >
    <button
      class="header-button"
      v-if="EventBus.isLoggedIn"
      @click="createNewGroup()"
    >
      New Group
    </button>
    <button
      class="header-button"
      v-if="EventBus.isLoggedIn"
      @click="getUsers()"
    >
      Users
    </button>
    <button
      class="header-button"
      v-if="EventBus.isLoggedIn"
      @click="getProfile()"
    >
      Profile
    </button>
    <button class="header-button" v-if="EventBus.isLoggedIn" @click="logout()">
      Logout
    </button>
    <button class="header-button" v-if="!EventBus.isLoggedIn" @click="login()">
      Login
    </button>
  </header>
  <RouterView />
  <footer class="footer">
    <i><b>Developed by @LucaMica02 for the Sapienza WASA course 2024</b></i>
  </footer>
</template>

<style>
.header {
  background-color: blanchedalmond;
  padding: 5px;
}

.header-button {
  font-size: 10px;
  margin: 5px;
}

.footer {
  background-color: blanchedalmond;
  padding: 10px;
  font-size: 5px;
  text-align: center;
}

.image {
  width: 25px;
  margin-right: 5px;
  border-radius: 50%;
  border: 1px solid black;
}

.clickable-container {
  cursor: pointer;
}
</style>
