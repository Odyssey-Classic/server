import React, { useEffect, useState } from 'react'
import MapCanvas from './components/MapCanvas'
import TilePalette from './components/TilePalette'
import { createMap as apiCreateMap, listMaps as apiListMaps, getMap as apiGetMap, updateMap as apiUpdateMap, type GameMap } from '../../api/maps'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function MapsModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()
    const [activeTile, setActiveTile] = useState<number>(0)
    const [maps, setMaps] = useState<GameMap[]>([])
    const [search, setSearch] = useState('')
    const [current, setCurrent] = useState<GameMap | null>(null)
    const [name, setName] = useState('')

    // Set title and menu items for this module
    useEffect(() => {
        setTitle('Maps')
        setMenuItems([
            { label: 'New Map', action: onNewMap },
            { label: 'Export All', action: () => console.log('Export all maps') },
            { label: 'Import', action: () => console.log('Import maps') }
        ])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Maps' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

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
        <div className="module-container">
            <main className="module-content" style={{ padding: 0 }}>
            </main>
        </div>
    )
}
