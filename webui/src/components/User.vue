<script>
export default {
  props: {
    user: Object, // Receive the user object as a prop
  },
  data() {
    return {
      conversation: [],
    };
  },
  methods: {
    // start a new chat with the user
    async startNewChat() {
      try {
        const userId = localStorage.getItem("authToken");
        const requestBody = [userId, this.user.resourceId].map((id) => ({
          resourceId: parseInt(id),
        }));
        console.log(requestBody);
        const response = await this.$axios.post(
          `/users/${userId}/conversations`,
          requestBody,
          {
            headers: { Authorization: userId },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Access token missing");
        } else if (response.status === 403) {
          alert("Not permitted");
        } else if (response.status === 409) {
          alert("Conversation already exists");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 201) {
          alert("Conversation created");
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },
  },
};
</script>

<template>
  <div class="user">
    <img :src="user.PhotoUrl" :alt="name" class="user-photo" />
    <h3>{{ "@" + user.username }}</h3>
    <button @click="startNewChat()" class="chat-button">💬</button>
  </div>
</template>

<style scoped>
.user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}

h3 {
  margin: 0;
}
</style>
