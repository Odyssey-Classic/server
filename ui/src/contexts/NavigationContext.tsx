import React, { createContext, useContext, ReactNode } from 'react'
import { type ModuleType } from '../components/Dashboard'

interface NavigationContextValue {
    navigateTo: (module: ModuleType) => void
}

const NavigationContext = createContext<NavigationContextValue | undefined>(undefined)

export function NavigationProvider({ children, navigateTo }: { children: ReactNode, navigateTo: (module: ModuleType) => void }) {
    return (
        <NavigationContext.Provider value={{ navigateTo }}>
            {children}
        </NavigationContext.Provider>
    )
}

export function useNavigation() {
    const context = useContext(NavigationContext)
    if (!context) {
        throw new Error('useNavigation must be used within a NavigationProvider')
    }
    return context
}
