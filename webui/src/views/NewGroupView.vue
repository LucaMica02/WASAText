<script setup>
import { RouterView } from "vue-router";
</script>

<script>
export default {
  data() {
    return {
      loggedUserId: null,
      groupName: "",
      groupDescription: "",
      usersToAdd: [],
      usersAdded: [],
    };
  },
  methods: {
    // Remove the user from 'usersToAdd' and add it to 'usersAdded'
    addMember(user) {
      this.usersToAdd = this.usersToAdd.filter(
        (u) => u.resourceId !== user.resourceId
      );
      this.usersAdded.push(user);
    },

    // Remove the user from 'usersAdded' and add it to 'usersToAdd'
    deleteMember(user) {
      if (user.resourceId != this.loggedUserId) {
        this.usersAdded = this.usersAdded.filter(
          (u) => u.resourceId !== user.resourceId
        );
        this.usersToAdd.push(user);
      }
    },

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
          this.usersToAdd = response.data;
          this.loggedUserId = localStorage.getItem("authToken");
          this.usersToAdd.forEach((user) => {
            user.resourceId == this.loggedUserId ? this.addMember(user) : null;
          });
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Create new group
    async createGroup() {
      if (this.groupName == "") {
        alert("Have to define a group name");
        return;
      }
      if (this.groupDescription == "") {
        alert("Have to define a group description");
        return;
      }
      try {
        const requestBody = {
          groupName: this.groupName,
          groupDescription: this.groupDescription,
          partecipants: this.usersAdded.map((user) => ({
            resourceId: user.resourceId,
          })),
        };
        const response = await this.$axios.post(
          `/users/${this.loggedUserId}/groups`,
          requestBody,
          {
            headers: { Authorization: this.loggedUserId },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Auth token missing");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 201) {
          alert("Group created");
          this.$router.replace("/");
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },
  },
  mounted() {
    this.fetchUsers();
  },
};
</script>

<template>
  <div class="profile-container">
    <div class="group-info">
      <!-- Textarea for set the group name -->
      <div class="change-group-name">
        <textarea
          v-model="groupName"
          placeholder="Set Group Name"
          rows="3"
          class="username-textarea"
        ></textarea>
      </div>

      <!-- Textarea for set the group description -->
      <div class="change-group-name">
        <textarea
          v-model="groupDescription"
          placeholder="Set Group Description"
          rows="3"
          class="username-textarea"
        ></textarea>
      </div>
    </div>

    <!-- show user added -->
    <div class="user-list">
      <h3>Users added:</h3>
      <div v-for="(user, index) in usersAdded" :key="index" class="user-item">
        <span @click="deleteMember(user)">
          {{ "@" + user.username }}
        </span>
      </div>
    </div>

    <!-- show user to add -->
    <div class="user-list">
      <h3>Select users to add:</h3>
      <div v-for="(user, index) in usersToAdd" :key="index" class="user-item">
        <span @click="addMember(user)">
          {{ "@" + user.username }}
        </span>
      </div>
    </div>
    <button class="header-button" @click="createGroup()">Create Group</button>
  </div>
  <main><RouterView /></main>
</template>

<style scoped>
.profile-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
  background-color: #f9f9f9;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.photo-section {
  position: relative;
}

.photo-container {
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
}

.profile-photo {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #4e73df;
}

.delete-photo-btn {
  background-color: #f04e4e;
  color: white;
  border: none;
  padding: 8px 20px;
  margin-top: 10px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s;
}

.delete-photo-btn:hover {
  background-color: #d03d3d;
}

.group-info {
  margin-top: 20px;
}

.group-info h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 20px;
}

.change-group-name,
.change-photo {
  margin-bottom: 10px;
}

.photo-upload-input {
  padding: 5px;
  border-radius: 5px;
  border: 1px solid #ddd;
  margin-bottom: 10px;
}

.btn {
  padding: 10px 20px;
  background-color: #4e73df;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn:hover {
  background-color: #2e59a6;
}
</style>
