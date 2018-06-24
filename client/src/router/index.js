import Vue from 'vue'
import Router from 'vue-router'
import List from '@/components/List'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: List
    },
    {
      path: '/list/:id',
      name: 'List',
      component: List
    }
  ]
})
