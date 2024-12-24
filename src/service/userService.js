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

export default {
  register,
  login,
  info,
};
