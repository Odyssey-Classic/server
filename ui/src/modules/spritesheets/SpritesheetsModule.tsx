import React, { useEffect, useState } from 'react'
import SpritesheetCanvas from './components/SpritesheetCanvas'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

interface SpritesheetData {
    imageUrl: string
    metadata: SpritesheetMetadata
}

interface SpritesheetMetadata {
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

export default function SpritesheetsModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()
    const [spritesheetData, setSpritesheetData] = useState<SpritesheetData | null>(null)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        setTitle('Spritesheets')
        setMenuItems([
            { label: 'Open Zip', action: handleOpenZip },
            { label: 'Clear', action: handleClear }
        ])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Spritesheets' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    const handleOpenZip = () => {
        const input = document.createElement('input')
        input.type = 'file'
        input.accept = '.zip'
        input.onchange = async (e) => {
            const file = (e.target as HTMLInputElement).files?.[0]
            if (!file) return

            try {
                const JSZip = (await import('jszip')).default
                const zip = await JSZip.loadAsync(file)

                // Look for spritesheet.png and spritesheet.json
                const pngFile = zip.file('spritesheet.png')
                const jsonFile = zip.file('spritesheet.json')

                if (!pngFile || !jsonFile) {
                    setError('ZIP file must contain both spritesheet.png and spritesheet.json')
                    return
                }

                // Extract the image
                const pngBlob = await pngFile.async('blob')
                const imageUrl = URL.createObjectURL(pngBlob)

                // Extract and parse the JSON
                const jsonText = await jsonFile.async('text')
                const metadata: SpritesheetMetadata = JSON.parse(jsonText)

                setSpritesheetData({ imageUrl, metadata })
                setError(null)
            } catch (err) {
                setError(`Failed to load spritesheet: ${err instanceof Error ? err.message : String(err)}`)
                console.error('Error loading spritesheet:', err)
            }
        }
        input.click()
    }

    const handleClear = () => {
        if (spritesheetData?.imageUrl) {
            URL.revokeObjectURL(spritesheetData.imageUrl)
        }
        setSpritesheetData(null)
        setError(null)
    }

    return (
        <div style={{ padding: '20px' }}>
            {error && (
                <div style={{
                    padding: '10px',
                    marginBottom: '20px',
                    backgroundColor: '#ff4444',
                    color: 'white',
                    borderRadius: '4px'
                }}>
                    {error}
                </div>
            )}

            {!spritesheetData && !error && (
                <div style={{
                    padding: '40px',
                    textAlign: 'center',
                    backgroundColor: '#2a2a2a',
                    borderRadius: '8px',
                    color: '#aaa'
                }}>
                    <p>No spritesheet loaded.</p>
                    <p>Click "Open Zip" to load a .zip file containing spritesheet.png and spritesheet.json</p>
                </div>
            )}

            {spritesheetData && (
                <div>
                    <div style={{ marginBottom: '20px', color: '#ccc' }}>
                        <h3>Spritesheet Info</h3>
                        <p>Total Frames: {Object.keys(spritesheetData.metadata.frames).length}</p>
                        {spritesheetData.metadata.meta?.size && (
                            <p>
                                Image Size: {spritesheetData.metadata.meta.size.w} x {spritesheetData.metadata.meta.size.h}
                            </p>
                        )}
                    </div>
                    <SpritesheetCanvas
                        imageUrl={spritesheetData.imageUrl}
                        metadata={spritesheetData.metadata}
                    />
                </div>
            )}
        </div>
    )
}
