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
        legal: Map<string, boolean>,
        handler: (index: number, event: React.MouseEvent<HTMLElement>) => void 
    }){
    return (
        <div className="board">
            {board && board.map((cell, index) => (
                <Cell key={index} piece={cell} index={index} valid={legal.has(index.toString())} handler={handler}/>
            ))}
        </div>
    )
}
