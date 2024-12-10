<script>
import Conversation from "../components/Conversation.vue";

export default {
  components: {
    Conversation,
  },
  data() {
    return {
      conversations: [],
      isLoading: true,
      selectedConversation: null,
      newMessage: "",
      usernames: {},
      currentUser: {},
    };
  },
  methods: {
    fetchCurrentUser() {
      this.currentUser["authToken"] = localStorage.getItem("authToken");
      this.currentUser["username"] = localStorage.getItem("username");
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

        this.conversations = [];
        for (const resource of response.data) {
          const conversation = await this.getConversation(
            resource["resourceId"]
          );
          conversation.resourceId = resource["resourceId"];
          this.conversations.push(conversation);
        }
      } catch (error) {
        console.error("Error: ", error);
      } finally {
        this.isLoading = false;
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
        return response.data;
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    async sendMessage(conversation) {
      try {
        const response = await this.$axios.post(
          `/users/${localStorage.getItem("authToken")}/conversations/${
            conversation.resourceId
          }/messages`,
          {
            repliedTo: 0,
            forwardedFrom: 0,
            type: "text",
            body: this.newMessage,
          },
          {
            headers: {
              Authorization: localStorage.getItem("authToken"),
            },
          }
        );
        this.newMessage = "";
        this.getConversations();
        const res = await this.getConversation(conversation.resourceId);
        this.selectConversation(res);
      } catch (error) {
        console.error("Error: ", error);
      }
    },

    // Select a conversation to view in detail
    selectConversation(conversation) {
      this.selectedConversation = conversation;
      this.getConversation(conversation.resourceId);
    },
  },
  mounted() {
    this.getConversations();
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

    <!-- This is where the selected conversation details will be displayed -->
    <div class="conversation-detail" v-if="selectedConversation">
      <div class="messages-container">
        <h2>{{ selectedConversation.conversationName }}</h2>
        <div v-if="selectedConversation.messages">
          <div
            v-for="(message, index) in selectedConversation.messages"
            :key="index"
          >
            <p>
              <strong
                >{{
                  usernames[message.sender] === currentUser["username"]
                    ? "You"
                    : usernames[message.sender]
                }}:</strong
              >
              {{ message.body }}
            </p>
          </div>
        </div>
        <div v-else>
          <p>No messages yet</p>
        </div>
      </div>
      <div class="message-input">
        <textarea
          v-model="newMessage"
          placeholder="Type the message.."
          rows="3"
          cols="40"
        ></textarea>
      </div>
      <button class="send-message" @click="sendMessage(selectedConversation)">
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

.conversations-list {
  width: 30%;
  border-right: 1px solid #ddd;
  padding-right: 10px;
}

.conversation-detail {
  width: 65%;
  padding-left: 10px;
}

.messages-container {
  max-height: 200px;
  overflow-y: auto;
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
</style>
