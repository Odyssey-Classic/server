import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        port: 5173
    },
    build: {
        // Output directly into the Go embed directory so the server can serve the UI without a copy step
        outDir: '../internal/web/dist',
        // Because outDir is outside the project root, explicitly allow cleaning it on build
        emptyOutDir: true,
        assetsDir: 'assets'
    }
})
