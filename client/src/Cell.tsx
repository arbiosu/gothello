/**
 * Represents a Cell on a Board
 */
export default function Cell({ piece, index, valid, handler }: { piece: string, index: number, valid: boolean, handler: (index: number, event: React.MouseEvent<HTMLElement>) => void }) {
    if (valid) {
        return (
            <div
                className="valid-cell"
                onClick={(event) => handler(index, event)}
            >{piece}
            </div>
        )
    }
    return (
        <div className="cell">{piece}</div>
    )
}