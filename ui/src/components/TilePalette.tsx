import React from 'react'

interface TilePaletteProps {
    activeTile: number
    onSelect: (id: number) => void
}

const SAMPLE_TILES = Array.from({ length: 20 }).map((_, i) => i)

export default function TilePalette({ activeTile, onSelect }: TilePaletteProps) {
    return (
        <div className="tile-palette">
            {SAMPLE_TILES.map((t) => (
                <div
                    key={t}
                    className={`tile ${t === activeTile ? 'active' : ''}`}
                    onClick={() => onSelect(t)}
                >
                    {t}
                </div>
            ))}
        </div>
    )
}
