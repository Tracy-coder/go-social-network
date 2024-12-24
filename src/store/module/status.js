import statusService from '@/service/statusService';

const statusModule = {
  namespaced: true,
  state: {
  },

  mutations: {
  },

  actions: {
    postStatus(context, message) {
      console.log(message);
      return new Promise((resolve, reject) => {
        console.log(message);
        statusService.postStatus(message).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
    fetchStatuses() {
      return new Promise((resolve, reject) => {
        statusService.fetchStatuses().then((res) => {
          console.log(res);
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
    fetchProfile() {
      return new Promise((resolve, reject) => {
        statusService.fetchProfile().then((res) => {
          console.log(res);
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
  },
};
export default statusModule;
