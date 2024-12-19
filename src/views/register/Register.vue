<template>
  <div class="register">
    <b-row class="mt-5">
      <b-col
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
              <b-card title="register">
          <b-form>
            <b-form-group label="username">
              <b-form-input
                v-model="$v.user.username.$model"
                type="text"
                required
                placeholder="Enter your name"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('username')">
                username lenth must between 2 and 12
              </b-form-invalid-feedback>
            </b-form-group>

            <b-form-group label="password">
              <b-form-input
                v-model="$v.user.password.$model"
                type="password"
                required
                placeholder="Enter your password"
                :state="validateState('password')"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('password')">
                password lenth must between 6 and 18
              </b-form-invalid-feedback>
            </b-form-group>

            <b-form-group label="e-mail">
              <b-form-input
                v-model="$v.user.email.$model"
                type="email"
                required
                placeholder="Enter your E-mail address"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('email')">
                please input valid E-mail address
              </b-form-invalid-feedback>
            </b-form-group>
            <b-form-group>
              <b-button
                @click="register"
                variant="outline-primary"
                block=""
              >register</b-button>
            </b-form-group>
          </b-form>
        </b-card>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators';

import customValidator from '@/helper/validator';
import { mapActions } from 'vuex';

export default {
  data() {
    return {
      user: {
        username: '',
        password: '',
        email: '',
      },
      validation: null,
    };
  },
  validations: {
    user: {
      username: {
        required,
        minLength: minLength(2),
        maxLength: maxLength(12),
      },
      email: {
        required,
        email: customValidator.emailValidator,
      },
      password: {
        required,
        minLength: minLength(6),
        maxLength: maxLength(18),
      },
    },
  },
  methods: {
    ...mapActions('userModule', { userRegister: 'register' }),
    validateState(name) {
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    register() {
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        console.log('reject');
        return;
      }
      console.log(this.user);
      this.userRegister(this.user).then(() => {
        this.$router.replace({ name: 'login' });
      }).catch((err) => {
        console.log(err.response.data.errMsg);
        this.$bvToast.toast(err.response.data.errMsg, {
          title: 'validation failed',
          variant: 'danger',
          solid: true,
        });
      });

      console.log('register');
    },
  },
};
</script>

<style lang="scss" scoped>
</style>
