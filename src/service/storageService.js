const PREFIX = 'jkdev_';

// user
const USER_PRIFIX = `${PREFIX}user_`;
const USER_TOKEN = `${USER_PRIFIX}token`;
const USER_INFO = `${USER_PRIFIX}info`;

const set = (key, data, expire) => {
  console.log(expire);
  const expiryDate = new Date(expire);
  const expiryTimestamp = expiryDate.getTime();

  const item = {
    value: data,
    expiry: expiryTimestamp,
  };

  localStorage.setItem(key, JSON.stringify(item));
};

const get = (key) => {
  const itemStr = localStorage.getItem(key);
  if (!itemStr) {
    return null;
  }
  const item = JSON.parse(itemStr);
  const now = new Date().getTime();
  console.log(now, item.expiry);
  if (item.expiry && now > item.expiry) {
    localStorage.removeItem(key);
    return null;
  }
  return item.value;
};
export default {
  set,
  get,
  USER_TOKEN,
  USER_INFO,
};
