import { app, store } from './app'

store.replaceState(global.__INITIAL_STATE__)

app.$mount("#app")