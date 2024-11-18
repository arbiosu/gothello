import { useState, useEffect, useRef } from 'react'
import Board from './Board'

interface Message {
    type: string;
    content: Content
}

interface Content {
    board: Array<string>;
    turn: string;
    legal: Map<number, boolean>;
    x: number;
    o: number;
    gameOver: boolean
}

const socketUrl: string = "ws://localhost:8081/ws"

export default function Game({ isActive }: { isActive: boolean }) {
    const [isConnected, setIsConnected] = useState<boolean>(false)
    const [board, setBoard] = useState<string[]>([])
    const [turn, setTurn] = useState<string>('X')
    const [legalMoves, setLegalMoves] = useState<Map<string, boolean>>(new Map<string, boolean>)
    const socketRef = useRef<WebSocket | null>(null)


    // socket event listeners
    useEffect(() => {
        socketRef.current = new WebSocket(socketUrl)
        socketRef.current.onopen = () => setIsConnected(true)
        socketRef.current.onclose = () => setIsConnected(false)
        socketRef.current.onmessage = (event: MessageEvent) => {
            const msg: Message = JSON.parse(event.data)
            if (msg.type === "gameState") {
                if (msg.content.gameOver) {
                    socketRef.current?.close()
                }
                // Set the GameData
                // TODO: bundle this all up?
                setBoard(() => {
                    return msg.content.board
                })
                setTurn(() => {
                    return msg.content.turn
                })
                setLegalMoves(() => {
                    const map: Map<string, boolean> = new Map(Object.entries(msg.content.legal))
                    return map
                })
            }
        }
        return () => {
            socketRef.current?.close()
        }
    }, [])
    // onClick handler
    const sendMove = (index: number, event: React.MouseEvent<HTMLElement>) => {
        event.preventDefault()
        socketRef.current?.send(JSON.stringify({type: "move", content: index}))
    }

    return (
        <>
            <Board board={board} handler={sendMove} legal={legalMoves} />
        </>
    )
}