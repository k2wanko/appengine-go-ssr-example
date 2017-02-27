import Vue from 'vue'
import Vuex from 'vuex'

import { mutations } from './mutations'
import { actions } from './actions'

Vue.use(Vuex)

const store = new Vuex.Store({
    state: {
        todo: {
            items: {}
        }
    },
    getters: {
        todoItems(state) {
            return Object.keys(state.todo.items).map(id => state.todo.items[id])
        }
    },
    actions,
    mutations
})

export default store