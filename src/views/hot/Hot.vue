<template>
  <div class="container mt-4">
    <div>
      <h2>Hot status</h2>
      <ul class="list-group">
        <li class="list-group-item" v-for="status in statuses" :key="status.ID">
          <status :entry="status" @unfollowUser="unfollowUser"
          @followUser="followUser"
          @toggleLikeStatus="toggleLikeStatus"></status>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import userService from '../../service/userService';
import Status from '../status/Status.vue';

export default {
  components: { Status },
  data() {
    return {
      newStatus: '',
      statuses: [
        // {
        //   ID: 1, userID: 'Alice', posted: '2024-12-24 10:30', message: '1',
        // },
        // {
        //   ID: 2, userID: 'Bob', posted: '2024-12-23 15:45', message: '2',
        // },
        // {
        //   ID: 3, userID: 'Charlie', posted: '2024-12-22 09:00', message: '3',
        // },
      ],
    };
  },

  created() {
    this.fetchHotStatus();
  },
  methods: {
    fetchHotStatus() {
      userService.fetchHotStatus().then((data) => {
        console.log(data);
        this.statuses = data.data.info;
        for (let i = 0; i < this.statuses.length; i += 1) {
          if (this.statuses[i].isLiked === true) {
            this.statuses[i].isLiked = true;
          } else {
            this.statuses[i].isLiked = false;
          }
        }
        console.log(this.statuses);
      })
        .catch((error) => {
          console.error('Error fetching statuses:', error);
        });
    },
    unfollowUser(userID) {
      userService.followAndUnfollow(userID, false).then((res) => {
        console.log(res);
        const entry = this.statuses.find((status) => status.userID === userID);
        if (entry) {
          this.$set(entry, 'isFollowed', false);
        }
      }).catch((err) => {
        console.log(err);
      });
    },
    followUser(userID) {
      userService.followAndUnfollow(userID, true).then((res) => {
        console.log(res);
        const entry = this.statuses.find((status) => status.userID === userID);
        if (entry) {
          this.$set(entry, 'isFollowed', true);
        }
      }).catch((err) => {
        console.log(err);
      });
    },
    toggleLikeStatus(ID, action) {
      console.log(ID);
      this.statuses = JSON.parse(JSON.stringify(this.statuses));
      userService.toggleLikeStatus(ID, action).then(() => {
        const entry = this.statuses.find((status) => status.ID === ID);
        if (entry) {
          this.$set(entry, 'isLiked', action);
        }
        console.log(this.statuses);
      }).catch((err) => {
        console.log(err);
      });
    },
  },

};
</script>
