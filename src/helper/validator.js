const telephoneValidator = (value) => /^1[3|4|5|6|7|8|9]\d{9}$/.test(value);
const emailValidator = (value) => /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/.test(value);
export default {
  telephoneValidator,
  emailValidator,
};
