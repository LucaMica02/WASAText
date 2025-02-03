<script setup>
import { RouterView } from "vue-router";
import User from "../components/User.vue";
</script>

<script>
export default {
  data() {
    return {
      users: [],
      searchQuery: "",
    };
  },
  methods: {
    // Get all the users
    async fetchUsers() {
      try {
        const response = await this.$axios.get(`/users`, {
          headers: { Authorization: localStorage.getItem("authToken") },
        });
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 404) {
          alert("User not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          response.data.forEach((user) => {
            user.PhotoUrl = this.getImagePath(user.PhotoUrl);
          });
          this.users = response.data;
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },
    getImagePath(PhotoUrl) {
      return this.$axios["defaults"]["baseURL"] + "/images?path=" + PhotoUrl;
    },
  },
  computed: {
    // filter user based on searchQuery
    filteredUsers() {
      return this.users.filter((user) =>
        user.username.toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
  },
  mounted() {
    this.fetchUsers();
  },
};
</script>

<template>
  <div class="search-bar">
    <input
      type="text"
      v-model="searchQuery"
      placeholder="Search a username"
      class="search-bar-input"
    />
  </div>
  <div v-if="filteredUsers.length > 0">
    <User v-for="user in filteredUsers" :key="user.id" :user="user" />
  </div>
  <div v-else>
    <p>Loading users...</p>
  </div>
  <main><RouterView /></main>
</template>

<style scoped>
.user-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.user {
  background-color: #f4f4f4;
  border-radius: 10px;
  margin: 10px;
  text-align: center;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s, box-shadow 0.2s;
  width: 80%;
  border: 2px solid #ddd;
}

.user:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.search-bar {
  display: inline;
  margin: 10px;
}

.search-bar-input {
  height: 20px;
  width: 200px;
}
</style>
