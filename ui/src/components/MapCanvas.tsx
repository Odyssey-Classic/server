import React, { useEffect, useRef } from 'react'
import * as PIXI from 'pixi.js'

interface MapCanvasProps {
    activeTile: number
}

const TILE_SIZE = 32
const MAP_DIM = 17

export default function MapCanvas({ activeTile }: MapCanvasProps) {
    const containerRef = useRef<HTMLDivElement | null>(null)
    const appRef = useRef<PIXI.Application | null>(null)

    useEffect(() => {
        if (!containerRef.current) return

        const app = new PIXI.Application({
            width: TILE_SIZE * MAP_DIM,
            height: TILE_SIZE * MAP_DIM,
            backgroundColor: 0x1e1e1e,
            resolution: window.devicePixelRatio || 1,
            antialias: true,
        })

        containerRef.current.appendChild(app.view as HTMLCanvasElement)
        appRef.current = app

        // simple mock: draw grid
        const grid = new PIXI.Container()
        app.stage.addChild(grid)

        for (let y = 0; y < MAP_DIM; y++) {
            for (let x = 0; x < MAP_DIM; x++) {
                const g = new PIXI.Graphics()
                g.lineStyle(1, 0x555555)
                g.beginFill(0x444444)
                g.drawRect(x * TILE_SIZE, y * TILE_SIZE, TILE_SIZE, TILE_SIZE)
                g.endFill()
                grid.addChild(g)
            }
        }

        return () => {
            app.destroy(true, { children: true })
            appRef.current = null
        }
    }, [])

    // TODO: hook up activeTile and click handlers to place tiles

    return <div ref={containerRef} style={{ width: TILE_SIZE * MAP_DIM, height: TILE_SIZE * MAP_DIM }} />
}
