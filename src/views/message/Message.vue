<template>
    <div>
      <b-card>
        <b-form-group>
        <b-form-input
        v-model="newMessage"
        type="text"
        required
        placeholder="Send a message"
      ></b-form-input>
      </b-form-group>
      <b-form-group>
      <b-button @click="sendMessage" variant="outline-primary">
      send</b-button>
      </b-form-group>
      </b-card>

      <div>
      <ul class="list-group">
        <li class="list-group-item" v-for="message in messageList" :key="message.ID">
        <strong>{{ message.senderName }}</strong> - {{ message.createdAt |formatDate }}
          <div>{{ message.content }}</div>
        </li>
      </ul>
    </div>
    </div>
</template>

<script>
import userService from '../../service/userService';

export default {
  data() {
    return {
      messageList: [],
      newMessage: '',
    };
  },
  created() {
    userService.getPendingMsg(this.$route.query.id).then((res) => {
      console.log(res);
      if (res.data && res.data.info && Array.isArray(res.data.info.info)
      && res.data.info.info.length > 0) {
        this.messageList = res.data.info.info.reverse();
        console.log(this.messageList);
      } else {
        console.warn('No messages or invalid response structure:', res.data);
        this.messageList = []; // 如果不符合条件，初始化为空数组
      }
    }).catch((error) => {
      console.error('Error fetching pending messages:', error);
      this.messageList = []; // 确保即使出错，messageList 也是空数组
    });
  },
  methods: {
    sendMessage() {
      console.log(this.newMessage);
      console.log(this.messageList);
      userService.postMessage(Number(this.$route.query.id), this.newMessage).then((res) => {
        console.log(res);
        this.messageList.unshift(res.data.info);
      });
      this.newMessage = '';
    },
  },
  filters: {
    formatDate(timestamp) {
      // console.log(timestamp);
      const date = new Date(timestamp);
      // console.log(date);
      return date.toLocaleString();
    },
  },
};
</script>
