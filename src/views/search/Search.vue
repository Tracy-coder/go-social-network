<template>

  <b-card>
    <b-card-text>
      <h3>search result:</h3>
          <UserList :userEntries="userEntries" @unfollowUser="unfollowUser"
          @followUser="followUser"></UserList>

    </b-card-text>
  </b-card>

</template>
<script>
import searchService from '../../service/searchService';
import userService from '../../service/userService';
import UserList from '../userList/UserList.vue';


export default {
  components: { UserList },
  data() {
    return {
      userEntries: [],
      searchQuery: '',
    };
  },

  created() {
    // 从 query 中获取数据
    this.searchQuery = this.$route.query.q;
    this.searchUser();
  },
  beforeRouteUpdate(to, from, next) {
    this.searchQuery = to.query.q;
    this.searchUser();
    next();
  },
  methods: {
    searchUser() {
      if (!this.searchQuery) return;
      console.log(this.searchQuery);
      searchService.searchUser(this.searchQuery).then((res) => {
        console.log(res);
        this.userEntries = res.data.userEntries;
        // for (let i = 0; i < this.userEntries.length; i += 1) {
        //   if (this.userEntries[i].isFollow === true) {
        //     this.userEntries[i].isFollow = true;
        //   } else {
        //     this.userEntries[i].isFollow = false;
        //   }
        // }
      }).catch((err) => {
        console.log(err);
      });
    },
    followUser(userID) {
      userService.followAndUnfollow(userID, true).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      const user = this.userEntries.find((entry) => entry.ID === userID);
      if (user) {
        this.$set(user, 'isFollow', true);
      }
      console.log(this.userEntries);
    },
    unfollowUser(userID) {
      userService.followAndUnfollow(userID, false).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      const user = this.userEntries.find((entry) => entry.ID === userID);
      if (user) {
        this.$set(user, 'isFollow', false);
      }
      console.log(this.userEntries);
    },
  },
};
</script>
