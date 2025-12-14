import React, { useState, useEffect } from 'react'
import TitleBar from './components/TitleBar'
import Breadcrumbs from './components/Breadcrumbs'
import Dashboard, { type ModuleType } from './components/Dashboard'
import { ServerInfoModule } from './modules/server-info'
import { UsersModule } from './modules/users'
import { AdministratorsModule } from './modules/administrators'
import { MapsModule } from './modules/maps'
import { ItemsModule } from './modules/items'
import { GuildsModule } from './modules/guilds'
import { ScriptsModule } from './modules/scripts'
import { TitleBarProvider } from './contexts/TitleBarContext'
import { NavigationProvider } from './contexts/NavigationContext'

function getModuleFromPath(): ModuleType {
    const path = window.location.pathname
    if (path.startsWith('/admin/server-info')) return 'server-info'
    if (path.startsWith('/admin/users')) return 'users'
    if (path.startsWith('/admin/administrators')) return 'administrators'
    if (path.startsWith('/admin/maps')) return 'maps'
    if (path.startsWith('/admin/items')) return 'items'
    if (path.startsWith('/admin/guilds')) return 'guilds'
    if (path.startsWith('/admin/scripts')) return 'scripts'
    return 'dashboard'
}

function getPathFromModule(module: ModuleType): string {
    if (module === 'dashboard') return '/admin'
    return `/admin/${module}`
}

export default function App() {
    const [currentModule, setCurrentModule] = useState<ModuleType>(getModuleFromPath())

    // Listen to browser back/forward navigation
    useEffect(() => {
        const handlePopState = () => {
            setCurrentModule(getModuleFromPath())
        }

        window.addEventListener('popstate', handlePopState)
        return () => window.removeEventListener('popstate', handlePopState)
    }, [])

    const navigateTo = (module: ModuleType) => {
        const path = getPathFromModule(module)
        window.history.pushState({ module }, '', path)
        setCurrentModule(module)
    }

    const renderModule = () => {
        switch (currentModule) {
            case 'server-info':
                return <ServerInfoModule />
            case 'users':
                return <UsersModule />
            case 'administrators':
                return <AdministratorsModule />
            case 'maps':
                return <MapsModule />
            case 'items':
                return <ItemsModule />
            case 'guilds':
                return <GuildsModule />
            case 'scripts':
                return <ScriptsModule />
            case 'dashboard':
            default:
                return <Dashboard onNavigate={navigateTo} />
        }
    }

    return (
        <TitleBarProvider>
            <NavigationProvider navigateTo={navigateTo}>
                <div className="admin-ui">
                    <TitleBar />
                    <Breadcrumbs />
                    <div className="admin-content">
                        {renderModule()}
                    </div>
                </div>
            </NavigationProvider>
        </TitleBarProvider>
    )
}
