import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Game from './Game'

function App() {
  const [newGame, setNewGame] = useState<boolean>(false)

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Gothello</h1>
      <div>
        <button onClick={() => (setNewGame(!newGame))}>New game</button>
      </div>
      {newGame && (
        <Game isActive={newGame} />
      )}

    </>
  )
}

export default App
