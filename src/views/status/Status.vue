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
    <div v-if="entry.getUrls && entry.getUrls.length">
<b-button v-show="showDelete" variant="primary" @click="$emit('deleteStatus',entry.ID)">
          deleteStatus
        </b-button>
  <div>
    <vue-gallery
      :images="entry.getUrls"
      :index="ind"
      @close="handleClose"
    ></vue-gallery>
    <div class="image-gallery">
      <img
        v-for="(url, i) in entry.getUrls"
        :key="i"
        :src="url"
        alt="Image"
        class="image-preview"
        @click="handleClick(i)"
        @error="handleImageError"
      />
    </div>
  </div>
    </div>
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
import VueGallery from 'vue-gallery';

export default {
  components: {
    VueGallery,
  },
  data() {
    return {
      ind: null,
    };
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
    showDelete() {
      return this.entry.userID === this.userInfo.ID;
    },
  },
  methods: {
    handleClick(i) {
      this.ind = i;
      console.log(this.ind);
    },
    handleImageError(event) {
      console.error('图片加载失败:', event.target.src);
    },
    handleClose() {
      console.log('关闭:', this.ind);
      this.ind = null;
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
<style scoped>
.image-gallery {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.image-preview {
  width: 100px;
  height: 100px;
  object-fit: cover;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.vue-gallery {
  background-color: rgba(0, 0, 0, 0.8);
}

</style>
