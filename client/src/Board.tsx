import Cell from './Cell'


/**
 * Board component
 */
export default function Board({ 
    board,
    legal,
    handler
    }:{ 
        board: Array<string>,
        legal: { [key: number]: boolean },
        handler: (index: number, event: React.MouseEvent<HTMLElement>) => void 
    }){
    return (
        <div className="board">
            {board && board.map((cell, index) => (
                <Cell key={index}
                piece={cell}
                index={index}
                valid={!!legal[index]}
                handler={handler}/>
            ))}
        </div>
    )
}
