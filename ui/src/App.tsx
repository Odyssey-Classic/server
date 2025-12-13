import React from 'react'
import TitleBar from './components/TitleBar'

export default function App() {
    return (
        <div className="admin-ui">
            <TitleBar title="Server Admin" />
            <div className="admin-content">
                <main className="dashboard">
                    <div className="welcome-message">
                        <h1>Welcome to Odyssey Server Admin</h1>
                        <p>Use the menu to navigate to different sections.</p>
                    </div>
                </main>
            </div>
        </div>
    )
}
