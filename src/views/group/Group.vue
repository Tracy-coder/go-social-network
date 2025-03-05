<template>

  <b-card>
    <b-card-text>

    <button @click="fetchFriends">new chat</button>
<!-- 这里写groupList -->
<div v-if="!showUserEntries">
  <ul class="list-group">
    <li class="list-group-item" v-for="chat in chatList" :key="chat.ID">
      <div @click="getPendingMsg(chat.ID)">
        <strong >{{ chat.ID }}</strong>
      </div>
          <b-button class="ml-auto" @click="leaveChat(chat.ID)" variant="outline-primary">
      Leave the chat
    </b-button>
    </li>
  </ul>
</div>
<div v-else>
<ul class="list-group">
  <li v-for="user in userEntries" :key="user.ID" class="list-group-item">
    <div>
      <input type="checkbox" :value="user.ID" v-model="picked" />
      <strong>ID:</strong> {{ user.ID }}
    </div>
    <p><strong>Username:</strong> {{ user.username }}</p>
  </li>
</ul>
<b-button @click="createChat"
variant="outline-primary" :disabled="picked.length === 0">done</b-button>
<b-button @click="cancel"
variant="outline-primary">cancel</b-button>
</div>
    </b-card-text>
  </b-card>

</template>
<script>
import userService from '../../service/userService';

export default {
  data() {
    return {
      chatList: [],
      showUserEntries: false,
      userEntries: [],
      picked: [],
    };
  },

  created() {
    this.fetchChatList();
  },
  methods: {
    fetchChatList() {
      userService.getChatList().then((res) => {
        console.log(res);
        this.chatList = res.data.info;
      }).catch((err) => {
        console.log(err);
      });
      this.showUserEntries = false;
    },
    cancel() {
      this.picked = [];
      this.showUserEntries = false;
    },
    fetchFriends() {
      userService.getFriends().then((res) => {
        console.log(res);
        this.userEntries = res.data.userEntries;
      }).catch((err) => {
        console.log(err);
      });
      this.showUserEntries = true;
    },
    createChat() {
      userService.createChat(this.picked).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      this.picked = [];
      this.fetchChatList();
    },
    getPendingMsg(id) {
      this.$router.push({ name: 'Message', query: { id } });
    },
    leaveChat(id) {
      userService.leaveChat(Number(id));
      console.log(Number(id));
      this.chatList = this.chatList.filter((chat) => chat.ID !== id);
    },
  },
};
</script>
