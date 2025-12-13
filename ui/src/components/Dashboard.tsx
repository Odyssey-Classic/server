import React from 'react'

export type ModuleType = 'dashboard' | 'server-info' | 'users' | 'maps' | 'items' | 'guilds' | 'scripts'

interface ModuleTileProps {
    title: string
    description: string
    icon: string
    onClick: () => void
}

function ModuleTile({ title, description, icon, onClick }: ModuleTileProps) {
    return (
        <div className="module-tile" onClick={onClick}>
            <div className="module-tile-icon">{icon}</div>
            <h2 className="module-tile-title">{title}</h2>
            <p className="module-tile-description">{description}</p>
        </div>
    )
}

interface DashboardProps {
    onNavigate: (module: ModuleType) => void
}

export default function Dashboard({ onNavigate }: DashboardProps) {
    const modules = [
        {
            id: 'server-info' as ModuleType,
            title: 'Server Info',
            description: 'View server status and statistics',
            icon: '🖥️'
        },
        {
            id: 'users' as ModuleType,
            title: 'Users',
            description: 'Manage user accounts and permissions',
            icon: '👥'
        },
        {
            id: 'maps' as ModuleType,
            title: 'Maps',
            description: 'Create and edit game maps',
            icon: '🗺️'
        },
        {
            id: 'items' as ModuleType,
            title: 'Items',
            description: 'Manage game items and inventory',
            icon: '⚔️'
        },
        {
            id: 'guilds' as ModuleType,
            title: 'Guilds',
            description: 'Configure guild settings',
            icon: '🛡️'
        },
        {
            id: 'scripts' as ModuleType,
            title: 'Scripts',
            description: 'Edit and manage game scripts',
            icon: '📜'
        }
    ]

    return (
        <div className="dashboard">
            <div className="dashboard-header">
                <h1>Odyssey Server Admin</h1>
                <p>Select a module to manage your server</p>
            </div>
            <div className="module-grid">
                {modules.map(module => (
                    <ModuleTile
                        key={module.id}
                        title={module.title}
                        description={module.description}
                        icon={module.icon}
                        onClick={() => onNavigate(module.id)}
                    />
                ))}
            </div>
        </div>
    )
}
