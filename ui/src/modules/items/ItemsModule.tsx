import React, { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function ItemsModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()

    useEffect(() => {
        setTitle('Items')
        setMenuItems([])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Items' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    return (
        <div className="module-container">
            <main className="module-content">
                <p>Item management and catalog will be displayed here.</p>
            </main>
        </div>
    )
}
