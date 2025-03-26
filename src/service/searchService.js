import request from '@/utils/request';

const searchUser = (message) => {
  console.log(message);
  return request.post('/api/v1/user/search', {
    expr: message,
  });
};

export default {
  searchUser,
};
