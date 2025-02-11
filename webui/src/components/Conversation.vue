<!-- Conversation.vue -->
<script>
export default {
  props: {
    conversation: Object, // Receive the conversation object as a prop
  },
};
</script>

<template>
  <div class="conversation-item">
    <div class="conversation-name">
      <strong>{{ conversation.conversationName }}</strong>
      <img
        :src="conversation.photoUrl"
        alt="Profile Photo"
        class="profile-photo"
      />
    </div>
    <div class="conversation-info">
      <span v-if="conversation.messages && conversation.messages.length > 0">
        Last message:
        {{ conversation.messages[conversation.messages.length - 1].body }}
        <div class="timestamp-prop">
          {{
            conversation.messages[conversation.messages.length - 1].timestamp
          }}
        </div>
      </span>
      <span v-else>No messages yet</span>
    </div>
  </div>
</template>

<style scoped>
.conversation-item {
  background-color: #f9f9f9;
  padding: 10px;
  border-radius: 8px;
  margin: 8px 0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.conversation-item:hover {
  background-color: #e8e8e8;
}

.conversation-name {
  font-size: 1.2rem;
  color: #333;
}

.conversation-info {
  font-size: 0.9rem;
  color: #777;
}

.profile-photo {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid #4e73df;
}

.timestamp-prop {
  margin-left: 5px;
  font-size: 10px;
}
</style>
