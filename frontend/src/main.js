import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import VueSimpleWebsocket from "vue-simple-websocket";
const app = createApp(App);
// app.use(VueSimpleWebsocket, "ws://localhost:8085/ws", {
//     reconnectEnabled: true,
//     reconnectInterval: 5000
// });
app.mount("#app");
// createApp(App).mount('#app')
