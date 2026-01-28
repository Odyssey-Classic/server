# Spritesheets Module

Admin module for managing and viewing game spritesheets.

## Features

### Current Implementation

1. **Local .zip File Loading**
   - Opens .zip files containing `spritesheet.png` and `spritesheet.json`
   - Validates that both required files are present
   - Extracts and parses spritesheet metadata

2. **PixiJS Rendering**
   - Displays the full spritesheet image
   - Renders each sprite individually at its correct position
   - Draws green borders around each sprite to visualize frame boundaries
   - Scrollable canvas for large spritesheets

3. **User Interface**
   - Clean tile-based dashboard entry point
   - "Open Zip" menu action to load spritesheets
   - "Clear" menu action to reset the view
   - Error handling with visual feedback
   - Info display showing frame count and image dimensions

## Usage

1. Navigate to the Spritesheets module from the admin dashboard
2. Click "Open Zip" in the menu bar
3. Select a .zip file containing:
   - `spritesheet.png` - the sprite atlas image
   - `spritesheet.json` - the spritesheet metadata in JSON format

The spritesheet will be rendered with each sprite positioned according to the metadata, with visible borders showing the sprite boundaries.

## Spritesheet JSON Format

The module expects spritesheet metadata in the following format:

```json
{
  "frames": {
    "sprite-name": {
      "frame": { "x": 0, "y": 0, "w": 32, "h": 32 },
      "rotated": false,
      "trimmed": false,
      "spriteSourceSize": { "x": 0, "y": 0, "w": 32, "h": 32 },
      "sourceSize": { "w": 32, "h": 32 }
    }
  },
  "meta": {
    "app": "texture-packer",
    "version": "1.0",
    "image": "spritesheet.png",
    "format": "RGBA8888",
    "size": { "w": 512, "h": 512 },
    "scale": "1"
  }
}
```

## Technical Details

### Components

- **SpritesheetsModule.tsx** - Main module component handling file loading and state
- **SpritesheetCanvas.tsx** - PixiJS canvas component for rendering sprites

### Dependencies

- `jszip` - For extracting .zip file contents
- `pixi.js` - For rendering sprites with WebGL/Canvas

### File Structure

```
spritesheets/
├── index.ts                     # Module export
├── SpritesheetsModule.tsx      # Main module component
└── components/
    └── SpritesheetCanvas.tsx   # PixiJS rendering component
```

## Future Enhancements

Potential additions for future development:

- Server-side API integration for saving/loading spritesheets
- Sprite selection and detail view
- Animation preview for sprite sequences
- Spritesheet validation and error reporting
- Export functionality
- Multi-spritesheet management
