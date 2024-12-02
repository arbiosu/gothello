export default function Sidebar() {
    return (
        <aside className='sidebar'>
            <h2>About Gothello</h2>
            <a href='https://github.com/arbiosu/gothello'>GitHub repository</a>
            <p>Gothello is a reimplementation of my <strong>Introduction to 
                Computer Science II final project.</strong>
                The final project spec included developing the board game
                <strong>Othello</strong> (AKA Reversi) in the terminal using Python.
                For the final project, I implemented the terminal game by creating a
                2D array to represent the board. 
                The '*' represent the borders of the board, 
                the '.' represent open spaces, and 'X' and 'O' represent the player's pieces. 
                The game was fully playable in the terminal, 
                allowing users to play locally with a friend. 
            </p>
            <p>As I continued my degree, I wanted to reimplement this project 
                using a new language and adding additional features. I aimed to make this game
                playable on the web against AI bots.
                <strong> First,</strong> I rewrote the terminal game in Go, optimizing the board 
                to be represented as a 100-element array instead of a 2D array.
                I added logic for a bot called Randy that only chooses
                random moves and then implemented a bot called Max that uses the
                <strong> minimax</strong> algorithm to choose moves. 
                Subsequently, I wrote the backend socket logic using the 
                gorilla/websocket library, and finally, 
                I developed the frontend in React using Vite.
            </p>
        </aside>
    )
}