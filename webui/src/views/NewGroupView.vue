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
        if (response.status === 200) {
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
        if (response.status === 201) {
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
      <h3>Group Name:</h3>
      <div class="change-group-name">
        <textarea
          v-model="groupName"
          placeholder="Set Group Name"
          rows="3"
          class="username-textarea"
        ></textarea>
      </div>

      <!-- Textarea for set the group description -->
      <h3>Group Description:</h3>
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

.user-item span {
  border: 2px solid #000;
  padding: 2px 5px;
  border-radius: 2px;
  cursor: pointer;
}

.user-item span:hover {
  border-color: #007bff;
  color: #007bff;
}

.header-button {
  background-color: #28a745;
  color: white;
  border: 5px;
  padding: 10px 20px;
  font-size: 16px;
  border-radius: 15px;
  cursor: pointer;
}

.header-button:hover {
  background-color: #218838;
  color: #fff;
  transform: scale(1.05);
  box-shadow: 0 4px 8px rgba(0, 123, 255, 0.3);
}
</style>
