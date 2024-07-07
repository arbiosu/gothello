// JS logic

document.getElementById("init").onclick = function(event) {
    let socket = new WebSocket("ws://localhost:7341/ws");
    console.log("Attempting WS connection...")
    let input = document.getElementById("name").value;
    document.getElementById("intro").remove();

    socket.onopen = () => {
        console.log("Successfully Connected to WS");
        data = {
            name: input,
            board: [],
            move: "", 
            legal: []
        };
        
        socket.send(JSON.stringify(data));
    };

    socket.onclose = (event) => {
        console.log("Server side closed the connection", event);
    };

    socket.onerror = (error) => {
        console.log("socket error: ", error);
    };

    socket.onmessage = (event) => {
        const game = JSON.parse(event.data);
        console.log(game);
        const board = game["board"];
        const moves = game["legal"];
        displayBoard(board);
        elems = displayLegalMoves(moves);
        for (let i = 0; i < elems.length; i++) {
            elems[i].addEventListener('click', function (e) {
                console.log(elems[i].id);
                let newData = {
                    name: input, 
                    board: board,
                    move: elems[i].id,
                    legal: []
                };
                console.log("Making move...");
                socket.send(JSON.stringify(newData));
            })
        }
    };
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

// legal: the legal moves for the player
function displayLegalMoves(legal) {
    const elems = [];
    for (let i = 0; i < legal.length; i++) {
        let move = document.getElementById(legal[i]);
        move.setAttribute("name", "legal");
        elems.push(move);
    }
    return elems
}