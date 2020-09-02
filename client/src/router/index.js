import Vue from 'vue'
import VueRouter from 'vue-router'
import Highlights from '../views/Highlights.vue'
import Highlight from '../views/Highlight.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'highlights',
    component: Highlights
  },
  {
    path: '/highlights/:id',
    name: 'highlight',
    component: Highlight
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
