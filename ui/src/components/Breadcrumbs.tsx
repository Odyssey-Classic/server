import React, { useEffect } from 'react'
import { useTitleBar } from '../contexts/TitleBarContext'

export default function Breadcrumbs() {
    const { breadcrumbs } = useTitleBar()

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
