import request from '@/utils/request';

const postStatus = (message, filenames) => {
  console.log(message, filenames);
  return request.post('/api/v1/user/post', { message, filenames });
};

const fetchStatuses = () => {
  return request.get('/api/v1/user/timeline');
};


export default {
  postStatus,
  fetchStatuses,
};
