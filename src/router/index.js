import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from '../views/Home.vue';
import userRoutes from './module/user';

Vue.use(VueRouter);

const routes = [
  {
    path: '/Home',
    name: 'Home',
    component: Home,
    meta: {
      auth: true,
    },
  },
  {
    path: '/search',
    name: 'Search',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/search/Search.vue'),
    meta: {
      auth: true,
    },
  },
  {
    path: '/group',
    name: 'Group',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/group/Group.vue'),
    meta: {
      auth: true,
    },
  },
  {
    path: '/message',
    name: 'Message',
    component: () => import('../views/message/Message.vue'),
    meta: {
      auth: true,
    },
  },
  {
    path: '/hot',
    name: 'Hot',
    component: () => import('../views/hot/Hot.vue'),
    meta: {
      auth: true,
    },
  },
  ...userRoutes,
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.meta.auth) {
    if (localStorage.getItem('jkdev_user_token')) {
      next();
    } else {
      router.push({ name: 'login' });
    }
  } else {
    next();
  }
});

export default router;
