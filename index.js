// JS logic

document.getElementById("init").onclick = function(event) {
    document.getElementById("init").remove();
    let socket = new WebSocket("ws://localhost:7339/ws");
    console.log("Attempting WS connection...")
    let input = document.getElementById("name").value;

    socket.onopen = () => {
        console.log("Successfully Connected to WS");
        data = {
            name: input,
            board: [],
            move: "", 
            legal: []
        };
        
        // socket.send(JSON.stringify({name:"Arbi"}));
        socket.send(JSON.stringify(data));
    };

    socket.onclose = (event) => {
        console.log("Client side closed the connection", event);
    };

    socket.onerror = (error) => {
        console.log("socket error: ", error);
    };

    socket.onmessage = (event) => {
        const game = JSON.parse(event.data);
        // const game = event.data.text().then(g=>console.log(g));
        console.log(game);
        const board = game["board"];
        const moves = game["legal"];
        displayBoard(board);
        elems = displayLegalMoves(moves);
    };
}

// Displays the current game state
// board: a 100 element array of strings representing the current game state
function displayBoard(board) {
    for (let i = 0; i < 100; i++) {
        if (i%10 === 0) {
            let row = document.createElement("div");
            document.getElementById("board").appendChild(row);
            row.className = "row";
        }
        let column = document.createElement("div");
        column.innerHTML = board[i];
        column.className = "column";
        let index = i.toString();
        column.id = index;
        row.appendChild(column);
    }
}

// legal: the legal moves for the player
function displayLegalMoves(legal) {
    const elems = [];
    for (let i = 0; i < legal.length; i++) {
        let move = document.getElementById(legal[i]);
        move.id = "legal";
        elems.push(move);
    }
    return elems
}