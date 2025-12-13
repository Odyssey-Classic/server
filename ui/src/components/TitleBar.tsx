import React, { useState } from 'react'

interface TitleBarProps {
    title?: string
}

export default function TitleBar({ title = 'Server Admin' }: TitleBarProps) {
    const [menuOpen, setMenuOpen] = useState(false)

    const toggleMenu = () => {
        setMenuOpen(!menuOpen)
    }

    return (
        <div className="title-bar">
            <div className="title-bar-left">
                <button 
                    className="icon-button hamburger-menu" 
                    onClick={toggleMenu}
                    aria-label="Menu"
                >
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <line x1="3" y1="6" x2="21" y2="6" />
                        <line x1="3" y1="12" x2="21" y2="12" />
                        <line x1="3" y1="18" x2="21" y2="18" />
                    </svg>
                </button>
                {menuOpen && (
                    <div className="dropdown-menu">
                        <div className="menu-item">File</div>
                        <div className="menu-item">Edit</div>
                        <div className="menu-item">View</div>
                        <div className="menu-item">Help</div>
                    </div>
                )}
            </div>

            <div className="title-bar-center">
                <h1 className="title-bar-title">{title}</h1>
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
