<template>
  <div>
    <b-navbar
      toggleable="lg"
      type="dark"
      variant="info"
    >
      <b-container>
        <b-navbar-brand @click="$router.push({name: 'Home'})">Go-social-network</b-navbar-brand>

        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse
          id="nav-collapse"
          is-nav
        >
                 <b-navbar-nav>
    <b-nav-item @click="$router.push({name:'Group'})">group</b-nav-item>
        <b-nav-item href="#">Disabled</b-nav-item>
      </b-navbar-nav>
          <b-navbar-nav class="ml-auto">
        <b-nav-form>
          <b-form-input size="sm" v-model="searchValue" class="mr-sm-2"
          placeholder="Search"></b-form-input>
          <b-button size="sm"  @click="search">Search</b-button>
        </b-nav-form>
            <b-nav-item-dropdown
             right
              v-if="userInfo"
            >
              <!-- Using 'button-content' slot -->
              <template v-slot:button-content>
                <em>{{userInfo.name}}</em>
              </template>
              <b-dropdown-item @click="$router.push({name:'profile'})">profile</b-dropdown-item>
              <b-dropdown-item @click="logout">logout</b-dropdown-item>
            </b-nav-item-dropdown>
            <div v-if="!userInfo">
              <b-nav-item
               v-if="$route.name != 'login'"
                @click="$router.replace({name: 'login'})"
              >login</b-nav-item>
              <b-nav-item
               v-if="$route.name != 'register'"
                @click="$router.replace({name: 'register'})"
              >register</b-nav-item>
            </div>
          </b-navbar-nav>
        </b-collapse>
      </b-container>
    </b-navbar>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {

  data() {
    return {
      searchValue: '',
    };
  },
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),

  methods: {
    ...mapActions('userModule', ['logout']),
    search() {
      console.log(this.searchValue);
      this.$router.push({ name: 'Search', query: { q: this.searchValue } });
    },
  },
};
</script>

<style scoped>
</style>
