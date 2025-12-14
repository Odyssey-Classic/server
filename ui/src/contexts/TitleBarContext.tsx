import React, { createContext, useContext, useState, ReactNode } from 'react'

export interface MenuItem {
    label: string
    action: () => void
}

export interface Breadcrumb {
    label: string
    onClick?: () => void
}

interface TitleBarContextValue {
    title: string
    setTitle: (title: string) => void
    menuItems: MenuItem[]
    setMenuItems: (items: MenuItem[]) => void
    breadcrumbs: Breadcrumb[]
    setBreadcrumbs: (crumbs: Breadcrumb[]) => void
}

const TitleBarContext = createContext<TitleBarContextValue | undefined>(undefined)

export function TitleBarProvider({ children }: { children: ReactNode }) {
    const [title, setTitle] = useState('Server Admin')
    const [menuItems, setMenuItems] = useState<MenuItem[]>([])
    const [breadcrumbs, setBreadcrumbs] = useState<Breadcrumb[]>([])

    return (
        <TitleBarContext.Provider value={{ title, setTitle, menuItems, setMenuItems, breadcrumbs, setBreadcrumbs }}>
            {children}
        </TitleBarContext.Provider>
    )
}

export function useTitleBar() {
    const context = useContext(TitleBarContext)
    if (!context) {
        throw new Error('useTitleBar must be used within a TitleBarProvider')
    }
    return context
}
