<script setup>
import { RouterView } from "vue-router";
</script>

<script>
export default {
  data() {
    return {
      groupConversation: {},
      newGroupName: "",
      newGroupDescription: "",
      newPhoto: null,
      showUserToAdd: false,
      users: {},
    };
  },
  methods: {
    // Get the user information
    fetchGroup() {
      this.groupConversation = JSON.parse(localStorage.getItem("group"));
      this.groupConversation.photoUrl = this.getImagePath();
      console.log(this.groupConversation);
    },

    // Show the user to add options
    openUserToAdd() {
      this.showUserToAdd = true;
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
          this.users = response.data;
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Update the group name
    async updateGroupName() {
      try {
        const authToken = localStorage.getItem("authToken");
        const response = await this.$axios.put(
          `/users/${authToken}/groups/${this.groupConversation.resourceId}/name`,
          { groupName: this.newGroupName },
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Invalid auth");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 404) {
          alert("Group not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          this.newGroupName = "";
          response.data.resourceId = this.groupConversation.resourceId;
          localStorage.setItem("group", JSON.stringify(response.data));
          this.fetchGroup();
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Update group description
    async updateGroupDescription() {
      try {
        const response = await this.$axios.put(
          `/users/${localStorage.getItem("authToken")}/groups/${
            this.groupConversation.resourceId
          }/description`,
          { groupDescription: this.newGroupDescription },
          {
            headers: { Authorization: localStorage.getItem("authToken") },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Invalid auth");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 404) {
          alert("Group not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          this.newGroupDescription = "";
          response.data.resourceId = this.groupConversation.resourceId;
          localStorage.setItem("group", JSON.stringify(response.data));
          this.fetchGroup();
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Update group photo
    async updatePhoto() {
      if (!this.newPhoto) {
        alert("Please select a photo");
        return;
      }
      const formData = new FormData();
      formData.append("photo", this.newPhoto);
      try {
        const response = await this.$axios.put(
          `/users/${localStorage.getItem("authToken")}/groups/${
            this.groupConversation.resourceId
          }/photo`,
          formData,
          {
            headers: {
              Authorization: localStorage.getItem("authToken"),
              "Content-Type": "multipart/form-data",
            },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Invalid auth");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 404) {
          alert("Group not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          this.newPhoto = null;
          response.data.resourceId = this.groupConversation.resourceId;
          localStorage.setItem("group", JSON.stringify(response.data));
          this.fetchGroup();
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Leave Group
    async leaveGroup() {
      try {
        const response = await this.$axios.delete(
          `/users/${localStorage.getItem("authToken")}/groups/${
            this.groupConversation.resourceId
          }/members`,
          {
            headers: { Authorization: localStorage.getItem("authToken") },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Invalid auth");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 404) {
          alert("Group not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 204) {
          alert("You have left the group");
          this.$router.replace("/");
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Add member
    async addMember(user) {
      try {
        const authToken = localStorage.getItem("authToken");
        console.log(authToken);
        const response = await this.$axios.put(
          `/users/${authToken}/groups/${this.groupConversation.resourceId}/members?userId=${user.resourceId}`,
          { body: "SenzaQuestoNonFunzionaNonSoPerchè" },
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Invalid auth");
        } else if (response.status === 403) {
          alert("Not authorized");
        } else if (response.status === 404) {
          alert("Group not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          alert("User already in the group");
        } else if (response.status === 201) {
          alert("Member Added Successfully");
          this.showUserToAdd = false;
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Handle file selection
    handleFileUpload(event) {
      const file = event.target.files[0];
      console.log(file.type);
      if (file && file.type.startsWith("image/")) {
        this.newPhoto = file;
      } else {
        alert("Please select a valid file.");
      }
    },

    // return the full image path
    getImagePath() {
      return (
        this.$axios["defaults"]["baseURL"] +
        "/images?path=" +
        this.groupConversation.photoUrl +
        "&t=" +
        new Date().getTime()
      );
    },
  },
  mounted() {
    this.fetchGroup();
    this.fetchUsers();
  },
};
</script>

<template>
  <div class="profile-container">
    <div class="photo-section">
      <div class="photo-container">
        <img
          :src="groupConversation.photoUrl"
          alt="Profile Photo"
          class="profile-photo"
        />
      </div>
    </div>
    <div class="group-info">
      <h2>{{ "@" + groupConversation.conversationName }}</h2>
      <!-- Textarea for changing the group name -->
      <div class="change-group-name">
        <textarea
          v-model="newGroupName"
          placeholder="Enter new group name"
          rows="3"
          class="username-textarea"
        ></textarea>
        <button @click="updateGroupName()" class="btn">
          Update Group Name
        </button>
      </div>

      <h2>{{ "Description: " + groupConversation.description }}</h2>
      <!-- Textarea for changing the group description -->
      <div class="change-group-name">
        <textarea
          v-model="newGroupDescription"
          placeholder="Enter new group description"
          rows="3"
          class="username-textarea"
        ></textarea>
        <button @click="updateGroupDescription()" class="btn">
          Update Group Description
        </button>
      </div>

      <!-- Photo upload form -->
      <div class="change-photo">
        <input
          type="file"
          accept="image/png, image/jpeg"
          @change="handleFileUpload"
          class="photo-upload-input"
        />
        <button @click="updatePhoto()" class="btn">Change Group Photo</button>
      </div>

      <button class="header-button" @click="leaveGroup()">Leave Group</button>
      <button class="header-button" @click="openUserToAdd()">Add Member</button>
    </div>

    <!-- show user list when -->
    <div v-if="showUserToAdd" class="user-list">
      <h3>Select a user to add:</h3>
      <div v-for="(user, index) in users" :key="index" class="user-item">
        <span @click="addMember(user)">
          {{ "@" + user.username }}
        </span>
      </div>
    </div>
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
