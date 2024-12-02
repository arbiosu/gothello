// websocket stuff
export interface Message {
    type: string;
    content: Content
}

export interface Content {
    board: Array<string>;
    turn: string;
    legal: { [key: number]: boolean };
    x: number;
    o: number;
    gameOver: boolean
}

const socketUrl: string = "ws://localhost:8081/ws"


type EventHandler = (...args: any[]) => void


class EventEmitter {
    private events: { [key: string]: EventHandler[] } = {}

    public on(event: string, handler: EventHandler) {
        if (!this.events[event]) {
            this.events[event] = []
        }
        this.events[event].push(handler)
    }

    public off(event: string, handler: EventHandler) {
        if (!this.events[event]) return
        this.events[event] = this.events[event].filter((h) => h !== handler)
    }

    public emit(event: string, ...args: any[]) {
        if (!this.events[event]) return
        this.events[event].forEach((handler) => handler(...args))
    }
}


class WebSocketClient extends EventEmitter {
    private static instance: WebSocketClient
    private client!: WebSocket
    private url: string

    private constructor(url: string) {
        super()
        this.url = url
        this.connect()
        window.addEventListener('beforeunload', this.handleBeforeUnload)
    }

    private connect() {
        this.client = new WebSocket(this.url)

        this.client.addEventListener("open", () => {
            this.emit("open")
        })

        this.client.addEventListener("message", (event) => {
            const msg: Message = JSON.parse(event.data)
            //console.log("Received msg: ", msg)
            this.emit("gameState", msg)
            if (msg.content.gameOver) {
                this.client.close()
                this.emit("close")
            }
        })

        this.client.addEventListener('close', () => {
            this.emit('close')
        })
    }

    public static getInstance(url: string): WebSocketClient {
        if (!WebSocketClient.instance) {
            WebSocketClient.instance = new WebSocketClient(url)
        }
        return WebSocketClient.instance
    }

    public send(data: any) {
        this.client.send(JSON.stringify(data))
    }

    public close() {
        this.client.close()
        window.removeEventListener('beforeunload', this.handleBeforeUnload)
    }

    private handleBeforeUnload() {
        this.close()
    }
}

export const socket = WebSocketClient.getInstance(socketUrl)
