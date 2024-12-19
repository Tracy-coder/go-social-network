<template>
  <div class="login">
    <b-row class="mt-5">
      <b-col
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
              <b-card title="login">
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
            <b-form-group>
              <b-button
                @click="login"
                variant="outline-primary"
                block=""
              >login</b-button>
            </b-form-group>
          </b-form>
        </b-card>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators';

import { mapActions } from 'vuex';

export default {
  data() {
    return {
      user: {
        username: '',
        password: '',
        mail: '',
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
      password: {
        required,
        minLength: minLength(6),
        maxLength: maxLength(18),
      },
    },
  },
  methods: {
    ...mapActions('userModule', { userlogin: 'login' }),
    validateState(name) {
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    login() {
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        console.log('reject');
        return;
      }

      this.userlogin(this.user).then(() => {
        this.$router.replace({ name: 'Home' });
      }).catch((err) => {
        console.log(err);
        this.$bvToast.toast(err.response.data.message, {
          title: 'validation failed',
          variant: 'danger',
          solid: true,
        });
      });
    },
  },
};
</script>

<style lang="scss" scoped>
</style>
