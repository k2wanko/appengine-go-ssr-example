import Vue from 'vue'
import App from './App.vue'

export default function createApp() {
  return new Vue({
    render: h => h(App)
  })
}