// JS logic
const newGameBtn = document.getElementById("newGame");

document.getElementById("init").onclick = function(event) {
    const socket = new WebSocket("ws://localhost:8081/ws");
    console.log(document.location.host);
    const name = document.getElementById("name").value;
    document.getElementById("intro").remove();

    socket.onopen = () => {
        console.log("Successfully Connected to WS");    
    };

    socket.onclose = (event) => {
        console.log("Server side closed the connection", event);
    };

    socket.onerror = (error) => {
        console.log("socket error: ", error);
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log(msg);
        if (msg.type === "gameState") {
            displayBoard(msg.content.board)
            displayScoreBoard(msg.content.o, msg.content.x)
            if (msg.content.turn === "X") {
                elems = displayLegalMoves(msg.content.legal)
                for (let i = 0; i < elems.length; i++) {
                    elems[i].addEventListener('click', function (e) {
                        console.log("Selected:", elems[i].id);
                        socket.send(JSON.stringify({type: "move", content: elems[i].id}));
                    })
                }
            }
        }
        if (msg.content.gameOver === true) {
            console.log("GameOver")
            socket.close()
        }

    };
    newGameBtn.addEventListener('click', () => {
        socket.send(JSON.stringify({type: 'newGame'}));
    });
}

// Displays the current game state
// board: a 100 element array of strings representing the current game state
function displayBoard(board) {
    let boardElem = document.getElementById("board");
    boardElem.innerHTML = "";
    let row;
    for (let i = 0; i < 100; i++) {
        if (i%10 === 0) {
            row = document.createElement("div");
            row.className = "row";
            boardElem.appendChild(row)
        }
        let column = document.createElement("div");
        column.innerHTML = board[i];
        column.className = "column";
        column.id = i.toString();
        row.appendChild(column);
    }
}

function displayScoreBoard(o, x) {
    let whiteScore = document.getElementById("white");
    let blackScore = document.getElementById("black");

    whiteScore.innerHTML = "O: " + o
    blackScore.innerHTML = "X: " + x
}

// legal: the legal moves for the player
function displayLegalMoves(legal) {
    const elems = [];
    // iterate over the keys of the legal moves map
    Object.keys(legal).forEach(entry => {
        let move = document.getElementById(entry);
        move.setAttribute("name", "legal");
        elems.push(move);
    });
    return elems
}