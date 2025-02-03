<script setup>
import { RouterView } from "vue-router";
</script>

<script>
export default {
  data() {
    return {
      currentUser: {},
      newUsername: "",
      newPhoto: null,
    };
  },
  methods: {
    // Get the user information
    async fetchUser() {
      try {
        const response = await this.$axios.get(
          `/users/${localStorage.getItem("authToken")}`,
          {
            headers: { Authorization: localStorage.getItem("authToken") },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 404) {
          alert("User not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          this.currentUser = response.data;
          this.currentUser["fullPhotoUrl"] = this.getImagePath();
          console.log(this.currentUser["fullPhotoUrl"]);
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },
    async updateUsername() {
      try {
        const response = await this.$axios.put(
          `/users/${localStorage.getItem("authToken")}/username`,
          { username: this.newUsername },
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
          alert("User not found");
        } else if (response.status === 409) {
          alert("Username not valid");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          localStorage.setItem("username", this.newUsername);
          this.newUsername = "";
          this.fetchUser();
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },
    async updatePhoto() {
      if (!this.newPhoto) {
        alert("Please select a photo");
        return;
      }
      const formData = new FormData();
      formData.append("photo", this.newPhoto);
      try {
        const response = await this.$axios.put(
          `/users/${localStorage.getItem("authToken")}/photo`,
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
          alert("User not found");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          this.newPhoto = null;
          this.currentUser["PhotoUrl"] = response.data["PhotoUrl"];
          this.currentUser["fullPhotoUrl"] = "";
          this.currentUser["fullPhotoUrl"] = this.getImagePath();
          this.fetchUser();
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
        this.currentUser.PhotoUrl +
        "&t=" +
        new Date().getTime()
      );
    },
  },
  mounted() {
    this.fetchUser();
  },
};
</script>

<template>
  <div class="profile-container">
    <div class="photo-section">
      <div class="photo-container">
        <img
          :src="currentUser.fullPhotoUrl"
          alt="Profile Photo"
          class="profile-photo"
        />
      </div>
    </div>
    <div class="user-info">
      <h2>{{ "@" + currentUser.username }}</h2>
      <!-- Textarea for changing the username -->
      <div class="change-username">
        <textarea
          v-model="newUsername"
          placeholder="Enter new username"
          rows="3"
          class="username-textarea"
        ></textarea>
        <button @click="updateUsername()" class="btn">Update Username</button>
      </div>
      <!-- Photo upload form -->
      <div class="change-photo">
        <input
          type="file"
          accept="image/png, image/jpeg"
          @change="handleFileUpload"
          class="photo-upload-input"
        />
        <button @click="updatePhoto()" class="btn">Change Photo</button>
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

.user-info {
  margin-top: 20px;
}

.user-info h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 20px;
}

.change-username,
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
