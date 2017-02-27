import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import Top from '../views/Top.vue'
import About from '../views/About.vue'
import Todo from '../views/Todo.vue'

export default new Router({
    mode: 'history',
    scrollBehavior: () => ({ y: 0 }),
    routes: [
        { path: '/', component: Top },
        { path: '/about', component: About },
        { path: '/todo', component: Todo },
    ],
})