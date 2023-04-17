<template>
    <div>
        <div v-for="(user, index) in users" :key="index" :style="{ backgroundColor: colors[index] }">
            {{ user }}
        </div>
    </div>
</template>
<script>
import { reactive, onMounted } from 'vue';
import ReconnectingWebSocket from 'reconnecting-websocket';

export default {
    setup() {
        const state = reactive({
            users: ['', '', '', '', ''],
            colors: ['red', 'green', 'blue', 'yellow', 'purple']
        });

        const updateUser = (goroutine, userID) => {
            console.log("ws")
            const index = Number(goroutine) - 1;
            state.users[index] = userID;
        };

        onMounted(() => {
            const socket = new ReconnectingWebSocket('ws://localhost:8085/ws');
            // socket.onmessage = (event) => {
            //     console.log("got message")
            //     const data = JSON.parse(event.data);
            //     const { goroutine, userID } = data;
            //     updateUser(goroutine, userID);
            // };
            socket.addEventListener('message', (event) => {
                console.log(event.data)
                    const data = JSON.parse(event.data);
                    const { goroutine, userID } = data;
                    updateUser(goroutine, userID);

            });
            socket.addEventListener('open', () => {
                console.log('WebSocket connected');
            });
        });

        return {
            ...state,
            updateUser
        };
    },
    // created () {
    //     this.$socketClient.onOpen = () => {
    //         console.log('socket connected')
    //     }
    //     this.$socketClient.onMessage = msg => {
    //         console.log("got message")
    //         console.log(JSON.parse(msg.data))
    //     }
    //     this.$socketClient.onClose = () => {
    //         console.log('socket closed')
    //     }
    //     this.$socketClient.onError = () => {
    //         console.log('socket error')
    //     }
    // }
}
</script>
