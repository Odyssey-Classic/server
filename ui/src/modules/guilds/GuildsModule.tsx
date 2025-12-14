import React, { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function GuildsModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()

    useEffect(() => {
        setTitle('Guilds')
        setMenuItems([])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Guilds' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    return (
        <div className="module-container">
            <main className="module-content">
                <p>Guild management will be displayed here.</p>
            </main>
        </div>
    )
}
