import React, { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function ServerInfoModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()

    useEffect(() => {
        setTitle('Server Info')
        setMenuItems([])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Server Info' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    return (
        <div className="module-container">
            <main className="module-content">
                <p>Server information and statistics will be displayed here.</p>
            </main>
        </div>
    )
}
