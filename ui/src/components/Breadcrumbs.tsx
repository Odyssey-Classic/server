import React, { useEffect } from 'react'
import { useTitleBar } from '../contexts/TitleBarContext'

export default function Breadcrumbs() {
    const { breadcrumbs } = useTitleBar()

    useEffect(() => {
        const content = document.querySelector('.admin-content') as HTMLElement
        if (content) {
            if (breadcrumbs.length > 0) {
                content.style.marginTop = 'calc(var(--title-bar-height) + var(--breadcrumbs-height))'
            } else {
                content.style.marginTop = 'var(--title-bar-height)'
            }
        }
    }, [breadcrumbs.length])

    if (breadcrumbs.length === 0) {
        return null
    }

    return (
        <div className="breadcrumbs">
            {breadcrumbs.map((crumb, index) => (
                <React.Fragment key={index}>
                    {index > 0 && <span className="breadcrumb-separator">/</span>}
                    {crumb.onClick ? (
                        <button
                            className="breadcrumb-link"
                            onClick={crumb.onClick}
                        >
                            {crumb.label}
                        </button>
                    ) : (
                        <span className="breadcrumb-current">{crumb.label}</span>
                    )}
                </React.Fragment>
            ))}
        </div>
    )
}
