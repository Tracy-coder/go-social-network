<template>

  <b-card>
    <b-card-text>
      <h3>hello <span v-if="userInfo">{{userInfo.username}}</span></h3>
    </b-card-text>
  <div class="user-info">
    <p><strong>Email:</strong> {{ userInfo.email }}</p>
    <p><strong>Following:</strong> {{ userInfo.following }}</p>
    <p><strong>Posts:</strong> {{ userInfo.posts }}</p>
    <p><strong>Signup:</strong> {{ userInfo.signup | formatDate}}</p>
  </div>
    <!-- <b-button href="#" variant="primary">edit profile</b-button> -->
    <div>
      <h2>Profile</h2>
      <ul class="list-group">
        <li class="list-group-item" v-for="status in statuses" :key="status.ID">
          <div>
            <strong>{{ status.username }}</strong> - {{ status.posted |formatDate }}
          </div>
          <div>{{ status.message }}</div>
        </li>
      </ul>
    </div>
  </b-card>

</template>
<script>
import { mapActions, mapState } from 'vuex';


export default {
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),

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
    this.fetchProfile();
  },
  methods: {
    ...mapActions('statusModule', { fetch: 'fetchProfile' }),
    fetchProfile() {
      this.fetch().then((data) => {
        console.log(data);
        this.statuses = data.data.info;
      })
        .catch((error) => {
          console.error('Error fetching statuses:', error);
        });
    },
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
