export function getTodoItems() {
    return fetch('/api/getTodoItems', { method: "POST" })
        .then(response => response.json())
}

export function addTodoItem(todo) {
    const body = new FormData()
    body.append('title', todo.title)
    return fetch('/api/addTodoItem', {
        method: 'POST',
        body,
    })
        .then(response => response.json())
}

export function removeTodoItem(id) {
    const body = new FormData()
    body.append('id', id)
    return fetch('/api/removeTodoItem', {
        method: 'POST',
        body,
    })
        .then(response => response.json())
}