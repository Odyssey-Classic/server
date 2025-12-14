# Title Bar Context

The title bar uses a React Context pattern that allows each module to control its own title dynamically.

## Adding a New Module

1. Create your module component
2. Import and use `useTitleBar` hook
3. Set title in `useEffect`
4. Add the module to App.tsx routing

No need to update any centralized configuration!
