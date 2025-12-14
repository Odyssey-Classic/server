# Title Bar Context

The title bar uses a React Context pattern that allows each module to control its own title and menu items dynamically.

## Usage in Modules

Each module can set its title and menu items using the `useTitleBar` hook:

```typescript
import { useEffect } from 'react'
import { useTitleBar } from '../../contexts/TitleBarContext'

export default function MyModule() {
    const { setTitle, setMenuItems } = useTitleBar()

    useEffect(() => {
        // Set the title that appears in the title bar
        setTitle('My Module')
        
        // Set menu items for the hamburger menu
        setMenuItems([
            { label: 'Action 1', action: () => handleAction1() },
            { label: 'Action 2', action: () => handleAction2() }
        ])
    }, [setTitle, setMenuItems])

    // ... rest of module
}
```

## Benefits

1. **Maintainable**: No centralized list of titles to update
2. **Flexible**: Each module controls its own title and menu
3. **Dynamic**: Titles and menus can change based on module state
4. **Type-safe**: Full TypeScript support

## API

### `useTitleBar()`

Returns an object with:
- `title: string` - Current title (read-only)
- `setTitle: (title: string) => void` - Set the title
- `menuItems: MenuItem[]` - Current menu items (read-only)
- `setMenuItems: (items: MenuItem[]) => void` - Set menu items

### `MenuItem`

```typescript
interface MenuItem {
    label: string
    action: () => void
}
```

## Adding a New Module

1. Create your module component
2. Import and use `useTitleBar` hook
3. Set title and menu items in `useEffect`
4. Add the module to App.tsx routing

No need to update any centralized configuration!
