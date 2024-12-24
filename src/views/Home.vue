<template>
  <div class="container mt-4">
    <div class="mb-3">
      <b-card>
        <b-form-group>
        <b-form-input
        v-model="newStatus"
        type="text"
        required
        placeholder="Post a status"
      ></b-form-input>
      </b-form-group>
      <b-form-group>
      <b-button @click="postStatus" variant="outline-primary">
      post status</b-button>
      </b-form-group>
      </b-card>
    </div>

    <div>
      <h2>Timeline</h2>
      <ul class="list-group">
        <li class="list-group-item" v-for="status in statuses" :key="status.ID">
          <div>
            <strong>{{ status.username }}</strong> - {{ status.posted |formatDate }}
          </div>
          <div>{{ status.message }}</div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
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
    this.fetchStatuses();
  },
  methods: {
    ...mapActions('statusModule', { post: 'postStatus', fetch: 'fetchStatuses' }),
    fetchStatuses() {
      this.fetch().then((data) => {
        console.log(data);
        this.statuses = data.data.info;
      })
        .catch((error) => {
          console.error('Error fetching statuses:', error);
        });
    },
    postStatus() {
      this.post(this.newStatus).then(() => {
        this.fetchStatuses();
      }).catch((err) => {
        console.log(err);
        this.$bvToast.toast(err.response.data.message, {
          title: 'post status failed',
          variant: 'danger',
          solid: true,
        });
      });
      this.newStatus = '';
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
