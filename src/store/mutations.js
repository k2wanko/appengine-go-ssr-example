import Vue from 'vue'

export const ADD_TODO_ITEM = "ADD_TODO_ITEM"
export const SET_TODO_ITEMS = "SET_TODO_ITEMS"
export const REMOVE_TODO_ITEM = "REMOVE_TODO_ITEM"

export const mutations = {
    [ADD_TODO_ITEM]: (state, todo) => {
        Vue.set(state.todo.items, todo.id, todo)
    },
    [SET_TODO_ITEMS]: (state, items) => {
        items.forEach(item => Vue.set(state.todo.items, item.id, item))
    },
    [REMOVE_TODO_ITEM]: (state, id) => {
        Vue.delete(state.todo.items, id)
    }
}