<template>

  <b-card>
    <b-card-text>
      <h3>hello <span v-if="userInfo">{{userInfo.username}}</span></h3>
    </b-card-text>
    <b-button href="#" variant="primary">edit profile</b-button>
  </b-card>
</template>
<script>
import { mapState } from 'vuex';

export default {
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),
};
</script>
<style scoped>
</style>
