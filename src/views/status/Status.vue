<template>
          <div>
            <strong>{{ entry.username }}</strong> - {{ entry.posted |formatDate }}
<b-button v-show="showFollow" variant="primary" @click="$emit('unfollowUser',entry.userID)">
          unfollow
        </b-button>
<b-button v-show="showUnfollow" variant="primary" @click="$emit('followUser',entry.userID)">
          follow
        </b-button>
          <div>{{ entry.message }}</div>
      <div>
        <b-button variant="primary" v-show="entry.isLiked"
        @click="$emit('toggleLikeStatus', entry.ID,false)">
          Unlike
        </b-button>
        <b-button variant="primary"
        v-show="!entry.isLiked" @click="$emit('toggleLikeStatus', entry.ID,true)">
          Like
        </b-button>
      </div>
          </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  data() {
    return {};
  },
  computed: {
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
      userStatuses: (state) => state.userStatus,
    }),
    showFollow() {
      return this.entry.isFollowed && this.entry.userID !== this.userInfo.ID;
    },
    showUnfollow() {
      return !this.entry.isFollowed && this.entry.userID !== this.userInfo.ID;
    },
  },
  props: {
    entry: Object,
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
