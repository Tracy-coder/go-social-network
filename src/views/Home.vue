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
  <b-button class="upload-button" @click="triggerUpload" variant="outline-primary">
      <b-icon icon="image"></b-icon>
    </b-button>
    <input
      type="file"
      ref="fileInput"
      style="display: none"
      accept="image/*,video/*"
      multiple
      @change="handleFileChange"
    />

    <b-container fluid class="mt-3">
      <b-row>
        <b-col
          v-for="(file, index) in previewFiles"
          :key="index"
          cols="4"
          class="mb-3"
        >
          <b-card
            class="preview-card"
            @click="test(file)"
          >
          <img :src="file.preview" alt="Image Preview" class="preview-image" />
            <b-button
              size="sm"
              variant="danger"
              class="remove-button"
              @click="removeFile(index)"
            >
              <b-icon icon="x"></b-icon>
            </b-button>
          </b-card>
        </b-col>
      </b-row>
    </b-container>
      <b-button @click="postStatus" variant="outline-primary">
      post status</b-button>
      </b-form-group>
      </b-card>
    </div>

    <div>
      <h2>Timeline</h2>
      <ul class="list-group">
        <li class="list-group-item" v-for="status in statuses" :key="status.ID">
          <status :entry="status" @unfollowUser="unfollowUser"
          @toggleLikeStatus="toggleLikeStatus"></status>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import axios from 'axios';
import userService from '../service/userService';
import Status from './status/Status.vue';

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
      selectedFiles: [],
      previewFiles: [],
    };
  },

  created() {
    this.fetchStatuses();
  },
  methods: {
    ...mapActions('statusModule', { post: 'postStatus', fetch: 'fetchStatuses' }),
    triggerUpload() {
      this.$refs.fileInput.click();
    },
    test(file) {
      console.log(file);
    },
    handleFileChange(event) {
      const { files } = event.target;
      this.previewFiles = [];
      this.selectedFiles = Array.from(files);

      this.selectedFiles.forEach((file) => {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.previewFiles.push({
            file,
            preview: e.target.result,
          });
        };
        reader.readAsDataURL(file);
      });
      console.log(this.previewFiles);
    },
    removeFile(index) {
      this.selectedFiles.splice(index, 1);
      this.previewFiles.splice(index, 1);
    },
    fetchStatuses() {
      this.fetch().then((data) => {
        console.log(data);
        this.statuses = data.data.info;
        // for (let i = 0; i < this.statuses.length; i += 1) {
        //   if (this.statuses[i].isLiked === true) {
        //     this.statuses[i].isLiked = true;
        //   } else {
        //     this.statuses[i].isLiked = false;
        //   }
        // }
        console.log(this.statuses);
      })
        .catch((error) => {
          console.error('Error fetching statuses:', error);
        });
    },
    postStatus() {
      let fileNames = null;

      if (this.selectedFiles.length > 0) {
        fileNames = this.selectedFiles.map((fileItem) => fileItem.name);
        console.log(fileNames);
      }

      // 发送请求获取上传URLs
      this.post({ message: this.newStatus, filenames: fileNames })
        .then(async (res) => {
          console.log(res.data.putUrls);
          if (res.data.length > 0) {
            const { putUrls } = res.data;

            const uploadPromises = putUrls.map((url, index) => {
              const file = this.selectedFiles[index];

              // 读取文件为ArrayBuffer
              return file.arrayBuffer().then((arrayBuffer) => {
                const blob = new Blob([arrayBuffer], { type: file.type });

                // 发送PUT请求上传文件
                return axios.put(url, blob, {
                  headers: {
                    'Content-Type': file.type,
                  },
                });
              });
            });

            try {
              await Promise.all(uploadPromises);
              console.log('所有文件上传成功');
              this.fetchStatuses();
              this.selectedFiles = [];
              this.previewFiles = [];
            } catch (err) {
              console.log(err);
            }
          }
        })
        .catch((err) => {
          console.log(err);
        });
      this.newStatus = '';
    },
    unfollowUser(userID) {
      userService.followAndUnfollow(userID, false).then((res) => {
        console.log(res);
      }).catch((err) => {
        console.log(err);
      });
      // this.fetchStatuses();
      this.statuses = this.statuses.filter((entry) => entry.userID !== userID);
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

<style scoped>
.preview-card {
  position: relative;
  overflow: hidden;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.preview-card img {
  width: 100%;
  height: auto;
  object-fit: cover;
}

.remove-button {
  position: absolute;
  top: 5px;
  right: 5px;
  background: red;
  color: white;
  border: none;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
