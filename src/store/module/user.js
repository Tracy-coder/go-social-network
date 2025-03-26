import storageService from '@/service/storageService';
import userService from '@/service/userService';

const userModule = {
  namespaced: true,
  state: {
    token: storageService.get(storageService.USER_TOKEN),
    userInfo: storageService.get(storageService.USER_INFO) ? JSON.parse(storageService.get(storageService.USER_INFO)) : null, //eslint-disable-line
    userStatus: [],
  },

  mutations: {
    SET_TOKEN(state, { token, expire }) {
      console.log(expire);
      storageService.set(storageService.USER_TOKEN, token, expire);
      state.token = token;
    },
    SET_USERINFO(state, userInfo) {
      storageService.set(storageService.USER_INFO, JSON.stringify(userInfo));
      state.userInfo = userInfo;
    },
    SET_USERSTATUS(state, userStatus) {
      state.userStatus = userStatus;
    },
  },

  actions: {
    register(context, { username, password, email }) {
      return new Promise((resolve, reject) => {
        userService.register({ username, password, email }).then(() => {
          resolve();
        }).catch((err) => {
          reject(err);
        });
      });
    },
    login(context, { username, password }) {
      return new Promise((resolve, reject) => {
        userService.login({ username, password }).then((res) => {
          console.log(res.data.expire);
          context.commit('SET_TOKEN', res.data);
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
    fetchProfile(context) {
      return Promise.all([
        userService.fetchProfile(),
        userService.info(),
      ])
        .then(([profileRes, userinfoRes]) => {
          console.log(profileRes.data.info);
          console.log(userinfoRes);
          context.commit('SET_USERSTATUS', profileRes.data.info);
          context.commit('SET_USERINFO', userinfoRes.data);
        })
        .catch((err) => {
          return err;
        });
    },
    deleteStatus(context, postID) {
      return new Promise((resolve, reject) => {
        userService.deleteStatus(postID)
          .then(() => {
            const newUserStatus = context.state.userStatus.filter((entry) => entry.ID !== postID);
            context.commit('SET_USERSTATUS', newUserStatus);
          })
          .then(() => {
            resolve();
          })
          .catch((err) => {
            reject(err);
          });
      });
    },
  },
};
export default userModule;
