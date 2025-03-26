<template>
              <div>
<li v-for="user in userEntries" :key="user.ID">
        <strong>ID:</strong>{{ user.ID }}
        <p><strong>Username:</strong>{{ user.username }}
          <button v-show="!user.isFollow" @click="$emit('followUser',user.ID)">
          follow
        </button>
<button v-show="user.isFollow" @click="$emit('unfollowUser',user.ID)">
          unfollow
        </button></p>

      </li>
          </div>
</template>

<script>

export default {
  data() {
    return {};
  },
  props: {
    userEntries: Array,
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
