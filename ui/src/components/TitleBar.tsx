import React from 'react'
import { useTitleBar } from '../contexts/TitleBarContext'

export default function TitleBar() {
    const { title } = useTitleBar()

    return (
        <div className="title-bar">
            <div className="title-bar-left">
            </div>

            <div className="title-bar-center">
                <h1 className="title-bar-title">
                    {title}
                </h1>
            </div>

            <div className="title-bar-right">
                <button className="icon-button" aria-label="Help">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <circle cx="12" cy="12" r="10" />
                        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
                        <line x1="12" y1="17" x2="12.01" y2="17" />
                    </svg>
                </button>
                <button className="icon-button" aria-label="Account">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                        <circle cx="12" cy="7" r="4" />
                    </svg>
                </button>
            </div>
        </div>
    )
}
