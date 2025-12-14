import React, { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function ScriptsModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()

    useEffect(() => {
        setTitle('Scripts')
        setMenuItems([])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Scripts' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    return (
        <div className="module-container">
            <main className="module-content">
                <p>Script management and editor will be displayed here.</p>
            </main>
        </div>
    )
}
