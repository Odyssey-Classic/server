import React, { useEffect, useRef } from 'react'
import * as PIXI from 'pixi.js'

interface SpritesheetCanvasProps {
    imageUrl: string
    metadata: {
        frames: {
            [key: string]: {
                frame: { x: number; y: number; w: number; h: number }
                rotated?: boolean
                trimmed?: boolean
                spriteSourceSize?: { x: number; y: number; w: number; h: number }
                sourceSize?: { w: number; h: number }
            }
        }
        meta?: {
            app?: string
            version?: string
            image?: string
            format?: string
            size?: { w: number; h: number }
            scale?: string
        }
    }
}

export default function SpritesheetCanvas({ imageUrl, metadata }: SpritesheetCanvasProps) {
    const containerRef = useRef<HTMLDivElement | null>(null)
    const appRef = useRef<PIXI.Application | null>(null)

    useEffect(() => {
        if (!containerRef.current) return

        let app: PIXI.Application | null = null

        const initPixi = async () => {
            // Determine canvas dimensions from metadata or use defaults
            const canvasWidth = metadata.meta?.size?.w || 800
            const canvasHeight = metadata.meta?.size?.h || 600

            // Create the PixiJS application
            app = new PIXI.Application()
            await app.init({
                width: canvasWidth,
                height: canvasHeight,
                backgroundColor: 0x1e1e1e,
                resolution: window.devicePixelRatio || 1,
                antialias: true,
            })

            if (!containerRef.current || !app) return
            containerRef.current.appendChild(app.canvas)
            appRef.current = app

            try {
                // Load the spritesheet image using Texture.from which handles blob URLs better
                const baseTexture = await PIXI.Assets.load({
                    src: imageUrl,
                    loadParser: 'loadTextures'
                }).catch(() => {
                    // Fallback to Texture.from for blob URLs
                    return PIXI.Texture.from(imageUrl)
                })

                console.log('Loaded texture:', baseTexture)

                if (!baseTexture) {
                    console.error('Failed to load texture from:', imageUrl)
                    return
                }

                // Wait for the texture to be ready (use source instead of baseTexture for v8)
                if (baseTexture.source && !baseTexture.source.valid) {
                    await new Promise(resolve => {
                        baseTexture.source.once('loaded', resolve)
                    })
                }

                console.log('Texture source:', baseTexture.source)
                console.log('Frame count:', Object.keys(metadata.frames).length)

                // Create a container to hold all sprites
                const spritesContainer = new PIXI.Container()
                app.stage.addChild(spritesContainer)

                // Iterate through each frame in the metadata
                Object.entries(metadata.frames).forEach(([frameName, frameData]) => {
                    const { x, y, w, h } = frameData.frame

                    console.log(`Creating sprite for ${frameName}: x=${x}, y=${y}, w=${w}, h=${h}`)

                    // Create a texture for this specific frame
                    const frameTexture = new PIXI.Texture({
                        source: baseTexture.source,
                        frame: new PIXI.Rectangle(x, y, w, h)
                    })

                    // Create a sprite from the frame texture
                    const sprite = new PIXI.Sprite(frameTexture)
                    sprite.position.set(x, y)

                    // Add the sprite to the container
                    spritesContainer.addChild(sprite)

                    // Draw a border around each sprite
                    const border = new PIXI.Graphics()
                    border.rect(x, y, w, h)
                    border.stroke({ width: 1, color: 0x00ff00, alpha: 0.5 })
                    spritesContainer.addChild(border)
                })

                console.log('Total sprites added:', spritesContainer.children.length)

            } catch (error) {
                console.error('Error loading spritesheet texture:', error)
            }
        }

        initPixi()

        return () => {
            if (app) {
                app.destroy(true, { children: true })
                appRef.current = null
            }
        }
    }, [imageUrl, metadata])

    return (
        <div
            ref={containerRef}
            style={{
                border: '1px solid #444',
                borderRadius: '4px',
                overflow: 'auto',
                maxWidth: '100%',
                maxHeight: '80vh'
            }}
        />
    )
}
