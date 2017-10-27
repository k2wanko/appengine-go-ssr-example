import Vue from 'vue'
import App from './App.vue'

export function createApp() {
  return new Vue({
    render: h => h(App)
  })
}
