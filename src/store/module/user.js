import storageService from '@/service/storageService';
import userService from '@/service/userService';

const userModule = {
  namespaced: true,
  state: {
    token: storageService.get(storageService.USER_TOKEN),
    userInfo: storageService.get(storageService.USER_INFO) ? JSON.parse(storageService.get(storageService.USER_INFO)) : null, //eslint-disable-line
  },

  mutations: {
    SET_TOKEN(state, token) {
      storageService.set(storageService.USER_TOKEN, token);
      state.token = token;
    },
    SET_USERINFO(state, userInfo) {
      storageService.set(storageService.USER_INFO, JSON.stringify(userInfo));
      state.userInfo = userInfo;
    },
  },

  actions: {
    register(context, { username, password, email }) {
      return new Promise((resolve, reject) => {
        userService.register({ username, password, email }).catch((err) => {
          reject(err);
        });
      });
    },
    login(context, { username, password }) {
      return new Promise((resolve, reject) => {
        userService.login({ username, password }).then((res) => {
          context.commit('SET_TOKEN', res.data.token);
          return userService.info();
        }).then((res) => {
          context.commit('SET_USERINFO', res.data);
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
    logout(context) {
      context.commit('SET_TOKEN', '');
      storageService.set(storageService.USER_TOKEN, '');
      context.commit('SET_USERINFO', '');
      storageService.set(storageService.USER_INFO, '');
      window.location.reload();
    },
  },
};
export default userModule;
