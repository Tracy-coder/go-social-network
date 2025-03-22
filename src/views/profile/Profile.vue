<template>
  <b-card>
    <b-card-text>
      <h3>hello <span v-if="userInfo">{{ userInfo.username }}</span></h3>
    </b-card-text>
    <div class="user-info">
      <p><strong>Email:</strong> {{ userInfo.email }}</p>
      <p @click="fetchFollowings"><strong>Following:</strong> {{ userInfo.followings }}</p>
      <p @click="fetchFollowers"><strong>Followed:</strong> {{ userInfo.followers }}</p>
      <p @click="fetchFriends"><strong>friends:</strong> {{ userInfo.friends }}</p>
      <p @click="fetchProfile"><strong>Posts:</strong> {{ userInfo.posts }}</p>
      <p><strong>Signup:</strong> {{ formattedSignup }}</p>
    </div>
<div v-if="!showUserEntries">
  <h2>Profile</h2>
  <ul class="list-group">
    <li class="list-group-item" v-for="status in userStatuses" :key="status.ID">
<status :entry="status" @deleteStatus="deleteStatus" @toggleLikeStatus="toggleLikeStatus"></status>
      <!-- <div>
        <strong>{{ status.username }}</strong> - {{ formatDate(status.posted) }}
      </div>
      <div>{{ status.message }}</div>
        <b-button variant="primary" @click="deleteStatus(status.ID)">
          Delete
        </b-button> -->
    </li>
  </ul>
</div>

<div v-else>
          <UserList :userEntries="userEntries" @unfollowUser="unfollowUser"
          @followUser="followUser"></UserList>
</div>
  </b-card>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import userService from '../../service/userService';
import UserList from '../userList/UserList.vue';
import Status from '../status/Status.vue';

export default {
  components: { UserList, Status },
  computed: {
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
      userStatuses: (state) => state.userStatus,
    }),
    formattedSignup() {
      return this.formatDate(this.userInfo.signup);
    },
  },
  data() {
    return {
      userEntries: [],
      showUserEntries: false,
    };
  },
  created() {
    this.fetchProfile();
  },
  methods: {
    ...mapActions('userModule', { fetchProfileAction: 'fetchProfile', deleteStatusAction: 'deleteStatus' }),
    fetchProfile() {
      this.fetchProfileAction().then((data) => {
        console.log(data);
      }).catch((error) => {
        console.error('Error getting profile:', error);
      });
      this.showUserEntries = false;
    },
    formatDate(timestamp) {
      console.log(timestamp);
      const date = new Date(timestamp);
      return date.toLocaleString();
    },
    fetchFollowings() {
      console.log('following');
      userService.getFollowings().then((res) => {
        console.log(res);
        this.userEntries = res.data.userEntries;
      }).catch((err) => {
        console.log(err);
      });
      this.showUserEntries = true;
    },
    fetchFollowers() {
      console.log('followers');
      userService.getFollowers().then((res) => {
        console.log(res);
        this.userEntries = res.data.userEntries;
      }).catch((err) => {
        console.log(err);
      });
      this.showUserEntries = true;
    },
    fetchFriends() {
      console.log('friends');
      userService.getFriends().then((res) => {
        console.log(res);
        this.userEntries = res.data.userEntries;
      }).catch((err) => {
        console.log(err);
      });
      this.showUserEntries = true;
    },
    followUser(userID) {
      console.log('follow');
      userService.followAndUnfollow(userID, true).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      const user = this.userEntries.find((entry) => entry.ID === userID);
      if (user) {
        this.$set(user, 'isFollow', true);
        // user.isFollow = true; // 将 isFollow 设置为 true
      }
      console.log(this.userEntries);
    },
    unfollowUser(userID) {
      console.log('unfollow');
      userService.followAndUnfollow(userID, false).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      const user = this.userEntries.find((entry) => entry.ID === userID);
      if (user) {
        this.$set(user, 'isFollow', false);
      }
    },
    deleteStatus(postID) {
      this.deleteStatusAction(postID).catch((err) => {
        console.log(err);
      });
    },
    toggleLikeStatus(ID, action) {
      console.log(ID);
      this.userStatuses = JSON.parse(JSON.stringify(this.userStatuses));
      userService.toggleLikeStatus(ID, action).then(() => {
        const entry = this.userStatuses.find((status) => status.ID === ID);
        if (entry) {
          this.$set(entry, 'isLiked', action);
        }
        console.log(this.userStatuses);
      }).catch((err) => {
        console.log(err);
      });
    },
  },
};
</script>
