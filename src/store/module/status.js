import statusService from '@/service/statusService';

const statusModule = {
  namespaced: true,
  state: {
  },

  mutations: {
  },

  actions: {
    postStatus(context, { message, filenames }) {
      console.log(message, filenames);
      return new Promise((resolve, reject) => {
        // console.log(message, files);
        statusService.postStatus(message, filenames).then((res) => {
          console.log(res);
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

  },
};
export default statusModule;
