import request from '@/utils/request';

const register = ({ email, password, username }) => {
  // console.log(username, password, email);
  return request.post('/api/v1/register', { username, email, password });
};

const login = ({ username, password }) => {
  return request.post('/api/v1/user/login', { username, password });
};

const info = () => {
  return request.get('/api/v1/user/info');
};

const followAndUnfollow = (otherID, action) => {
  return request.post('/api/v1/user/follow', { otherID, action });
};

const toggleLikeStatus = (ID, action) => {
  return request.post('/api/v1/user/like', { ID, action });
};

const fetchProfile = () => {
  return request.get('/api/v1/user/profile');
};

const getFollowings = () => {
  return request.get('/api/v1/user/followings');
};

const getFollowers = () => {
  return request.get('/api/v1/user/followers');
};

const getFriends = () => {
  return request.get('/api/v1/user/friends');
};

const getChatList = () => {
  return request.get('/api/v1/user/chatlist');
};

const createChat = (picked) => {
  return request.post('/api/v1/user/chats', { memberID: picked });
};

const getPendingMsg = (ID) => {
  console.log(ID);
  return request.get('/api/v1/user/chat', { params: { ID } });
};

const postMessage = (ID, message) => {
  return request.post('/api/v1/user/chat', { ID, message });
};

// 定义一个函数，用于退出聊天
const leaveChat = (ID) => {
  // 发送一个delete请求，请求路径为/api/v1/user/chat，参数为ID
  return request.delete('/api/v1/user/chat', { params: { ID } });
};

const fetchHotStatus = () => {
  return request.get('/api/v1/user/hot');
};

const deleteStatus = (postID) => {
  return request.delete('/api/v1/user/post', { params: { postID } });
};
export default {
  register,
  login,
  info,
  followAndUnfollow,
  fetchProfile,
  getFollowings,
  getFollowers,
  getFriends,
  getChatList,
  createChat,
  getPendingMsg,
  postMessage,
  leaveChat,
  toggleLikeStatus,
  fetchHotStatus,
  deleteStatus,
};
