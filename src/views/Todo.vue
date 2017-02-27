<template>
    <div>
        <h1>Todo</h1>
        <input v-model="newTodo" @keyup.enter="addTodo(newTodo)" placeholder="New ToDo">
        <ul>
            <li v-for="item in todoItems">
                <item :item="item" @removeClick="removeTodoItem"></item>
            </li>
        </ul>
    </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Item from '../Item.vue'

import { getTodoItems, addTodoItem, removeTodoItem } from '../store/actions'

function fetchItems(store) {
    return getTodoItems().then(items => {
        store.commit(SET_TODO_ITEMS, items)
    })
}

export default {
    data () {
        return {
            newTodo: "",
        }
    },
    beforeMount () {
        // this.getTodoItems()
    },
    preFetch (store) {
        
    },
    computed: {
        ...mapGetters([
            'todoItems'
        ])
    },
    methods: {
        addTodo (newTodo) {
            const todo = {
                id: this.todoItems.length + 1,
                title: newTodo
            }
            this.addTodoItem(todo)
            this.newTodo = ""

        },
        ...mapActions({
            getTodoItems,
            addTodoItem,
            removeTodoItem
        })
    },
    components: {
        Item
    }
}
</script>

<style scoped>
li {
  display: block;
  margin: 0 10px;
}
</style>