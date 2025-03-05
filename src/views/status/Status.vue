<template>
              <div>
            <strong>{{ entry.username }}</strong> - {{ entry.posted |formatDate }}
<button  @click="$emit('unfollowUser',entry.userID)">
          unfollow
        </button>
          </div>
</template>

<script>


export default {
  data() {
    return {};
  },
  props: {
    entry: Object,
  },
  filters: {
    formatDate(timestamp) {
      // console.log(timestamp);
      const date = new Date(timestamp / 1000000);
      // console.log(date);
      return date.toLocaleString();
    },
  },

};
</script>
