import React, { useState } from 'react'
import MapCanvas from './components/MapCanvas'
import TilePalette from './components/TilePalette'

export default function App() {
    const [activeTile, setActiveTile] = useState<number>(0)

    return (
        <div className="admin-ui">
            <header className="toolbar">
                <h1>Odyssey Admin UI</h1>
                <div className="actions">
                    <button>Save</button>
                </div>
            </header>
            <main className="editor">
                <aside className="palette">
                    <TilePalette activeTile={activeTile} onSelect={setActiveTile} />
                </aside>
                <section className="canvas">
                    <MapCanvas activeTile={activeTile} />
                </section>
                <aside className="inspector">
                    <h3>Inspector</h3>
                </aside>
            </main>
        </div>
    )
}
