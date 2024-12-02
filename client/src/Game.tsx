import { useState, useEffect } from 'react'
import { socket, Content, Message } from "../socket.ts"
import Board from './Board.tsx'
import Sidebar from './Sidebar.tsx'


export default function Game() {
    const [newGame, setNewGame] = useState<boolean>(false)
    const [gameState, setGameState] = useState<Content | null>(null)

    const handleNewGame = () => {
        setNewGame(!newGame)
        socket.send({ type: "init" })
    }

    useEffect(() => {
        const handleMessage = (msg: Message) => {
            setGameState(msg.content)
            console.log("set game state")
        }

        const handleClose = () => {
            console.log('Websocket Connection Closed')
        }

        socket.on('gameState', handleMessage)
        socket.on('close', handleClose)

        return () => {
            socket.off('gameState', handleMessage)
            socket.off('close', handleClose)
            //socket.close()
        }
    }, [])


    const onClickHandler = (index: number, event: React.MouseEvent<HTMLElement>) => {
        socket.send({ type: "move", content: index })
        console.log("Send event: ", event)
        console.log("Sent move: ", index)
    }

    return (
        <div className='app-container'>
            <Sidebar />
            <div className='content'>
                {gameState ? (
                    <main>
                        <p>Current Turn: {gameState.turn}</p>
                        <p>Score: X - {gameState.x}  O - {gameState.o}</p>
                        <Board board={gameState.board} legal={gameState.legal} handler={onClickHandler} />
                        <button onClick={handleNewGame}>New Game</button>

                    </main>
                ) : (
                    <main>
                        <p>Waiting...</p>
                        <button onClick={handleNewGame}>New Game</button>
                    </main>
                )}
            </div>
        </div>
    )
}