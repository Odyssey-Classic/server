import React, { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'
import { useNavigation } from '../../contexts/NavigationContext'

export default function UsersModule() {
    const { setTitle, setMenuItems, setBreadcrumbs } = useTitleBar()
    const { navigateTo } = useNavigation()

    useEffect(() => {
        setTitle('Users')
        setMenuItems([])
        setBreadcrumbs([
            { label: 'Home', onClick: () => navigateTo('dashboard') },
            { label: 'Users' }
        ])
    }, [setTitle, setMenuItems, setBreadcrumbs, navigateTo])

    return (
        <div className="module-container">
            <main className="module-content">
                <p>User management will be displayed here.</p>
            </main>
        </div>
    )
}
