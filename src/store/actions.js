import * as api from '../api'
import { SET_TODO_ITEMS, ADD_TODO_ITEM, REMOVE_TODO_ITEM } from './mutations'

export const getTodoItems = "getTodoItems"
export const addTodoItem = "addTodoItem"
export const removeTodoItem = "removeTodoItem"

export const actions = {
    [getTodoItems]: ({ commit }) => {
        return api.getTodoItems().then(items => commit(SET_TODO_ITEMS, items))
    },
    [addTodoItem]: ({ commit }, todo) => {
        return api.addTodoItem(todo).then(todo => commit(ADD_TODO_ITEM, todo))
    },
    [removeTodoItem]: ({ commit }, id) => {
        return api.removeTodoItem(id).then(todo => commit(REMOVE_TODO_ITEM, todo.id))
    }
}