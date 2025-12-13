import React, { useEffect, useState } from 'react'
import MapCanvas from './components/MapCanvas'
import TilePalette from './components/TilePalette'
import { createMap as apiCreateMap, listMaps as apiListMaps, getMap as apiGetMap, updateMap as apiUpdateMap, type GameMap } from '../../api/maps'

export default function MapsModule() {
    const [activeTile, setActiveTile] = useState<number>(0)
    const [maps, setMaps] = useState<GameMap[]>([])
    const [search, setSearch] = useState('')
    const [current, setCurrent] = useState<GameMap | null>(null)
    const [name, setName] = useState('')

    // load list on mount and whenever search changes (debounced via useEffect)
    useEffect(() => {
        let cancelled = false
        const t = setTimeout(async () => {
            try {
                const data = await apiListMaps(search || undefined)
                if (!cancelled) setMaps(data)
            } catch (e) {
                // eslint-disable-next-line no-console
                console.error('list maps failed', e)
            }
        }, 200)
        return () => {
            cancelled = true
            clearTimeout(t)
        }
    }, [search])

    // new blank map (local state only until saved)
    const onNewMap = () => {
        const blank: GameMap = {
            id: 0,
            name: 'Untitled Map',
            tags: [],
            attributes: {},
            last_updated: new Date().toISOString(),
            version: 1,
            tiles: [],
            links: {},
        }
        setCurrent(blank)
        setName(blank.name)
    }

    const onSave = async () => {
        try {
            if (!current || !name.trim()) return
            if (current.id === 0) {
                const saved = await apiCreateMap(name.trim())
                setCurrent(saved)
                setName(saved.name)
                // refresh list
                const data = await apiListMaps(search || undefined)
                setMaps(data)
            } else {
                const updated: GameMap = { ...current, name: name.trim() }
                await apiUpdateMap(updated)
                setCurrent(updated)
                const data = await apiListMaps(search || undefined)
                setMaps(data)
            }
        } catch (e) {
            // eslint-disable-next-line no-console
            console.error('save failed', e)
        }
    }

    const onLoad = async (id: number) => {
        try {
            const m = await apiGetMap(id)
            setCurrent(m)
            setName(m.name)
        } catch (e) {
            // eslint-disable-next-line no-console
            console.error('load failed', e)
        }
    }

    return (
        <>
            <header className="toolbar" style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <h1 style={{ marginRight: 'auto' }}>Map Editor</h1>
                <button onClick={onNewMap}>New Map</button>
                <input
                    placeholder="Map name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    style={{ width: 240 }}
                />
                <button onClick={onSave} disabled={!name.trim()}>
                    Save
                </button>
            </header>
            <main className="editor" style={{ display: 'grid', gridTemplateColumns: '280px 1fr 280px', gap: 12 }}>
                <aside className="list" style={{ padding: 8, borderRight: '1px solid #333' }}>
                    <h3 style={{ marginTop: 0 }}>Maps</h3>
                    <input
                        placeholder="Search by name"
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        style={{ width: '100%', marginBottom: 8 }}
                    />
                    <div style={{ maxHeight: 500, overflow: 'auto' }}>
                        {maps.length === 0 && <div style={{ opacity: 0.7 }}>No maps</div>}
                        {maps.map((m) => (
                            <div
                                key={m.id}
                                role="button"
                                onClick={() => onLoad(m.id)}
                                style={{
                                    padding: '6px 8px',
                                    cursor: 'pointer',
                                    background: current?.id === m.id ? '#2a2a2a' : 'transparent',
                                }}
                            >
                                {m.name} {m.id ? <span style={{ opacity: 0.6 }}>#{m.id}</span> : null}
                            </div>
                        ))}
                    </div>
                </aside>
                <section className="canvas">
                    <MapCanvas activeTile={activeTile} />
                </section>
                <aside className="palette">
                    <TilePalette activeTile={activeTile} onSelect={setActiveTile} />
                </aside>
            </main>
        </>
    )
}
