<script>
import Conversation from "../components/Conversation.vue";

export default {
  components: {
    Conversation,
  },
  data() {
    return {
      showForwardOptions: false,
      forwardMessage: null,
      conversations: [],
      isLoading: true,
      selectedConversation: null,
      newMessage: "",
      usernames: {},
      currentUser: {},
      users: [],
      replyTo: 0,
      newPhoto: null,
    };
  },
  methods: {
    fetchCurrentUser() {
      this.currentUser["authToken"] = localStorage.getItem("authToken");
      this.currentUser["username"] = localStorage.getItem("username");
    },
    openForwardOptions(message) {
      this.forwardMessage = message;
      this.showForwardOptions = true;
    },
    // Load the user data
    async getUsers() {
      try {
        const response = await this.$axios.get(`/users`, {
          headers: {
            Authorization: localStorage.getItem("authToken"),
          },
        });
        const users = response.data;
        this.users = response.data;
        for (let user of users) {
          this.usernames[user.resourceId] = user.username;
        }
      } catch (error) {
        console.error("Error: " + error);
      }
    },
    // Get all the user conversations
    async getConversations() {
      try {
        const response = await this.$axios.get(
          `/users/${localStorage.getItem("authToken")}/conversations`,
          {
            headers: {
              Authorization: localStorage.getItem("authToken"),
            },
          }
        );

        const localConversations = [];
        for (const resource of response.data) {
          const conversation = await this.getConversation(
            resource["resourceId"]
          );
          conversation.resourceId = resource["resourceId"];
          conversation.photoUrl = this.getImagePath(conversation.photoUrl);
          localConversations.push(conversation);
        }
        localConversations.sort((a, b) => {
          const lastMessageA = a.messages
            ? new Date(a.messages[a.messages.length - 1].timestamp)
            : new Date(0);
          const lastMessageB = b.messages
            ? new Date(b.messages[b.messages.length - 1].timestamp)
            : new Date(0);
          return lastMessageB - lastMessageA;
        });
        this.conversations = localConversations;
        this.selectedConversation = await this.getConversation(
          this.selectedConversation.resourceId
        );
      } catch (error) {
        console.error("Error: ", error);
      } finally {
        this.isLoading = false;
      }
    },

    async addReceiver(conversationId, messageId) {
      try {
        const userId = localStorage.getItem("authToken");
        const response = await this.$axios.put(
          `/conversations/${conversationId}/messages/${messageId}/receivers/${userId}`
        );
      } catch (error) {
        // console.error("Error: ", error);
      }
    },

    async addReader(conversationId, messageId) {
      try {
        const userId = localStorage.getItem("authToken");
        const response = await this.$axios.put(
          `/conversations/${conversationId}/messages/${messageId}/readers/${userId}`
        );
      } catch (error) {
        // console.error("Error: ", error);
      }
    },

    async addReaders(conversation) {
      if (conversation.messages) {
        for (const m of conversation.messages) {
          await this.addReader(conversation.resourceId, m.resourceId);
        }
      }
    },

    // Get a specific user conversation
    async getConversation(conversationId) {
      try {
        const response = await this.$axios.get(
          `/users/${localStorage.getItem(
            "authToken"
          )}/conversations/${conversationId}`,
          {
            headers: {
              Authorization: localStorage.getItem("authToken"),
            },
          }
        );
        response.data.resourceId = conversationId;
        if (response.data.messages) {
          for (const m of response.data.messages) {
            await this.addReceiver(conversationId, m.resourceId);
          }
        }
        return response.data;
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Send a message
    async sendMessage(conversation, path) {
      if (!this.newPhoto && this.newMessage === "") {
        alert("Cannot send empty message");
        return;
      }
      try {
        const response = await this.$axios.post(
          `/users/${localStorage.getItem("authToken")}/conversations/${
            conversation.resourceId
          }/messages`,
          {
            repliedTo: this.replyTo,
            forwardedFrom: 0,
            type: path === "" ? "text" : "image",
            body: path === "" ? this.newMessage : path,
          },
          {
            headers: {
              Authorization: localStorage.getItem("authToken"),
            },
          }
        );
        this.newMessage = "";
        this.replyTo = 0;
        this.newPhoto = null;
        this.getConversations();
        const res = await this.getConversation(conversation.resourceId);
        this.selectConversation(res);
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Send photo
    async sendPhoto(conversation) {
      // save image
      if (!this.newPhoto) {
        alert("Please select a photo");
        return;
      }
      const formData = new FormData();
      const path =
        localStorage.getItem("username") + "_" + new Date().getTime();
      formData.append("photo", this.newPhoto);
      try {
        const response = await this.$axios.post(
          `/images/newImage?path=${path}`,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        if (response.status === 200) {
          // send the message
          this.sendMessage(conversation, response.data.path);
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Forward a message
    async forwardMessageToConversation(conversation) {
      try {
        const authToken = localStorage.getItem("authToken");
        const requestBody = { resourceId: conversation.resourceId };
        const response = await this.$axios.post(
          `/users/${authToken}/conversations/${this.selectedConversation.resourceId}/messages/${this.forwardMessage.resourceId}/forward`,
          requestBody,
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Access token missing");
        } else if (response.status === 403) {
          alert("Not permitted");
        } else if (response.status === 404) {
          alert("Message missing");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 201) {
          alert("Forwarded");
          let messageCopy = { ...this.forwardMessage };
          messageCopy.sender = authToken;
          if (conversation.messages === null) {
            conversation.messages = [messageCopy];
          } else {
            conversation.messages.push(messageCopy);
          }
          this.showForwardOptions = false;
          this.forwardMessage = null;
          this.selectedConversation = await this.getConversation(
            conversation.resourceId
          );
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Delete a message
    async deleteMessage(index) {
      try {
        const authToken = localStorage.getItem("authToken");
        const response = await this.$axios.delete(
          `/users/${authToken}/conversations/${this.selectedConversation.resourceId}/messages/${this.selectedConversation.messages[index].resourceId}`,
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Access token missing");
        } else if (response.status === 403) {
          alert("Not permitted");
        } else if (response.status === 404) {
          alert("Message missing");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 204) {
          this.selectedConversation.messages.splice(index, 1);
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    async uncommentMessage(message) {
      try {
        const authToken = localStorage.getItem("authToken");
        const response = await this.$axios.delete(
          `/users/${authToken}/conversations/${this.selectedConversation.resourceId}/messages/${message.resourceId}/comment`,
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Access token missing");
        } else if (response.status === 404) {
          alert("Resource missing");
        } else if (response.status === 500) {
          alert("Server Error");
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Comment a message
    async commentMessage(message) {
      try {
        const requestBody = { emoji: "LIKE" };
        const authToken = localStorage.getItem("authToken");
        const response = await this.$axios.put(
          `/users/${authToken}/conversations/${this.selectedConversation.resourceId}/messages/${message.resourceId}/comment`,
          requestBody,
          {
            headers: { Authorization: authToken },
          }
        );
        if (response.status === 400) {
          alert("Bad Request");
        } else if (response.status === 401) {
          alert("Access token missing");
        } else if (response.status === 403) {
          alert("Not permitted");
        } else if (response.status === 404) {
          alert("Resource missing");
        } else if (response.status === 500) {
          alert("Server Error");
        } else if (response.status === 200) {
          // delete the comment
          this.uncommentMessage(message);
        }
      } catch (error) {
        console.error("Error: ", error);
      } finally {
        // update the view
        this.getConversations();
        const res = await this.getConversation(
          this.selectedConversation.resourceId
        );
        this.selectConversation(res);
      }
    },

    // start a new chat with the user
    async startNewChat(user) {
      try {
        const userId = localStorage.getItem("authToken");
        const requestBody = [userId, user.resourceId].map((id) => ({
          resourceId: parseInt(id),
        }));
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
          return response.data;
        }
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // If not exists create the conversation then forward the message
    async forwardMessageToUserConversation(user) {
      var conversation = this.conversations.find(
        (conv) => conv.conversationName === user.username
      );
      if (conversation) {
        this.forwardMessageToConversation(conversation);
      } else {
        conversation = await this.startNewChat(user);
        this.forwardMessageToConversation(conversation);
        this.conversations.push(conversation);
      }
    },

    // Manage send button
    manageSend(conversation) {
      this.newPhoto
        ? this.sendPhoto(conversation)
        : this.sendMessage(conversation, "");
    },

    // Handle file selection
    handleFileUpload(event) {
      const file = event.target.files[0];
      if (file && file.type.startsWith("image/")) {
        this.newPhoto = file;
      } else {
        alert("Please select a valid file.");
      }
    },

    // Select message to reply
    selectMessageToReply(message) {
      this.replyTo =
        this.replyTo === message.resourceId ? 0 : message.resourceId;
    },

    // Select a conversation to view in detail
    selectConversation(conversation) {
      this.replyTo = 0;
      this.newPhoto = null;
      this.selectedConversation = conversation;
      this.getConversation(conversation.resourceId);
      this.addReaders(conversation);
    },

    // Return the message content
    getOriginalMessage(messageId) {
      const message = this.selectedConversation.messages.find(
        (m) => m.resourceId == messageId
      );
      return message;
    },

    // Go in the group info page
    goToGroupInfo() {
      localStorage.setItem("group", JSON.stringify(this.selectedConversation));
      this.$router.replace("/groupInfo");
    },

    // Return full image path
    getImagePath(path) {
      return this.$axios["defaults"]["baseURL"] + "/images?path=" + path;
    },

    startPolling() {
      this.getConversations(); // Ottieni subito le conversazioni iniziali

      // Imposta il polling ogni 5000ms (5 secondi)
      this.pollingInterval = setInterval(() => {
        this.getConversations();
      }, 2000);
    },
  },
  mounted() {
    this.startPolling();
    this.fetchCurrentUser();
    this.getUsers();
  },
};
</script>

<template>
  <div class="main-container">
    <!-- This is where all the conversations will be displayed -->
    <div class="conversations-list">
      <div v-if="isLoading">Loading..</div>
      <div v-if="!isLoading && conversations.length === 0">
        No Conversation Yet
      </div>
      <ul v-if="!isLoading && conversations.length > 0">
        <li
          v-for="(conversation, index) in conversations"
          :key="index"
          @click="selectConversation(conversation)"
        >
          <Conversation :conversation="conversation" />
        </li>
      </ul>
    </div>

    <!-- show conversation list when click on forward -->
    <div v-if="showForwardOptions" class="forward-options">
      <h3>Select a conversation to forward the message:</h3>
      <h4>Groups:</h4>
      <div
        v-for="(conversation, index) in conversations.filter(
          (c) => !c.isPrivate
        )"
        :key="index"
        class="conversation-item"
        @click="forwardMessageToConversation(conversation)"
      >
        {{ conversation.conversationName }}
      </div>
      <h4>Users:</h4>
      <div
        v-for="(user, index) in users.filter(
          (u) => u.resourceId != currentUser.authToken
        )"
        :key="index"
        class="conversation-item"
        @click="forwardMessageToUserConversation(user)"
      >
        {{ user.username }}
      </div>
    </div>

    <!-- This is where the selected conversation details will be displayed -->
    <div class="conversation-detail" v-if="selectedConversation">
      <div class="messages-container">
        <h2 class="conversation-title">
          {{ selectedConversation.conversationName }}
        </h2>
        <button
          class="header-button"
          v-if="!selectedConversation.isPrivate"
          @click="goToGroupInfo()"
        >
          GroupInfo
        </button>
        <div v-if="selectedConversation.messages">
          <div
            v-for="(message, index) in selectedConversation.messages"
            :key="index"
          >
            <div :class="{ 'reply-message': message.repliedTo !== 0 }">
              <span v-if="message.repliedTo !== 0"
                ><strong>Replied to:</strong>
                <span v-if="!getOriginalMessage(message.repliedTo)"
                  >Message Cancelled</span
                >
                <span
                  v-else-if="
                    getOriginalMessage(message.repliedTo).type === 'text'
                  "
                  >{{ getOriginalMessage(message.repliedTo).body }}</span
                >
                <span v-else
                  ><img
                    :src="
                      this.$axios['defaults']['baseURL'] +
                      '/images?path=' +
                      getOriginalMessage(message.repliedTo).body
                    "
                    alt="Profile Photo"
                    class="message-photo"
                /></span>
              </span>
              <p>
                <span v-if="message.status === 'received'"> ✔️ </span>
                <span v-if="message.status === 'read'"> ✔️✔️ </span>
                <strong
                  >{{
                    usernames[message.sender] === currentUser["username"]
                      ? "You"
                      : usernames[message.sender]
                  }}:</strong
                >
                <span class="message-body" v-if="message.type === 'text'">
                  {{ message.body }}
                </span>
                <span class="message-body" v-else>
                  <img
                    :src="
                      this.$axios['defaults']['baseURL'] +
                      '/images?path=' +
                      message.body
                    "
                    alt="Profile Photo"
                    class="message-photo"
                  />
                </span>

                <!-- if is user's message the user can delete the message -->
                <span
                  v-if="usernames[message.sender] === currentUser['username']"
                  class="delete-icon"
                  @click="deleteMessage(index)"
                >
                  &#10006;
                </span>
                <span class="forward-icon" @click="openForwardOptions(message)">
                  &#8594;
                </span>
                <button
                  class="header-button"
                  @click="selectMessageToReply(message)"
                >
                  reply
                </button>
                <span v-if="message.resourceId === replyTo"> &#x25C9; </span>
                <span class="comment-icon" @click="commentMessage(message)">
                  &#x2764;{{ message.comments }}
                </span>
                <span v-if="message.forwardedFrom !== 0" class="forwarded-mark">
                  <b><i>*Forwarded</i></b>
                </span>
                <span class="timestamp-prop"> {{ message.timestamp }} </span>
              </p>
            </div>
          </div>
        </div>
        <div v-else>
          <p>No messages yet</p>
        </div>
      </div>
      <!-- Text Area -->
      <div class="message-input">
        <textarea
          v-model="newMessage"
          placeholder="Type the message.."
          rows="3"
          cols="40"
        ></textarea>
      </div>
      <!-- Image Area -->
      <div class="change-photo">
        <input
          type="file"
          accept="image/png, image/jpeg"
          @change="handleFileUpload"
          class="photo-upload-input"
        />
      </div>
      <button class="send-message" @click="manageSend(selectedConversation)">
        Send
      </button>
    </div>
  </div>

  <main><RouterView /></main>
</template>

<style scoped>
.main-container {
  display: flex;
  justify-content: space-between;
  padding: 20px;
}

.conversation-title {
  display: inline;
}

.group-info-title {
  font-size: 12px;
  color: darkgreen;
}

.group-info-title:hover {
  color: rgb(1, 50, 1);
}

.conversations-list {
  width: 30%;
  border-right: 1px solid #ddd;
  padding-right: 10px;
}

.conversation-detail {
  width: 80%;
  padding-left: 10px;
}

.messages-container {
  overflow-y: auto;
  max-height: 800px;
  border: 1px solid #ddd;
}

.conversations-list ul {
  list-style-type: none;
  padding: 0;
}

.conversations-list li {
  cursor: pointer;
  margin-bottom: 10px;
}

.delete-icon {
  cursor: pointer;
  color: red;
  margin-left: 10px;
  font-size: 20px;
}

.delete-icon:hover {
  color: darkred;
}

.forward-icon {
  cursor: pointer;
  color: blue;
  margin-left: 10px;
  font-size: 20px;
}

.forward-icon:hover {
  color: darkblue;
}

.forward-options h3 {
  font-size: 16px;
  margin-bottom: 10px;
  color: #613b3b;
  font-weight: normal;
}

.conversation-item {
  padding: 8px;
  margin: 4px 0;
  cursor: pointer;
  font-size: 14px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: #f0f0f0;
}

.comment-icon {
  cursor: pointer;
  color: black;
  margin-left: 10px;
  font-size: 20px;
}

.comment-icon:hover {
  color: darkgray;
}

.reply-message {
  margin-left: 20px;
  background-color: #f0f0f0;
  border-radius: 5px;
  padding: 5px;
}

.forwarded-mark {
  font-size: 10px;
  margin-left: 10px;
  color: darkred;
}

.timestamp-prop {
  margin-left: 5px;
  font-size: 10px;
  display: block;
}

.message-photo {
  width: 120px;
  height: 120px;
  object-fit: cover;
  border: 3px solid #4e73df;
  margin: 5px;
}

.message-body {
  margin-left: 10px;
}
</style>
