import request from '@/utils/request';

const postStatus = (message) => {
  // console.log(username, password, email);
  return request.post('/api/v1/user/post', { message });
};

const fetchStatuses = () => {
  return request.get('/api/v1/user/timeline');
};

const fetchProfile = () => {
  return request.get('/api/v1/user/profile');
};

export default {
  postStatus,
  fetchProfile,
  fetchStatuses,
};
