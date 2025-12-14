# Data Package

This package provides an abstraction layer for the server's data directory structure.

## Overview

The `data` package defines a `Root` interface that provides access to various subdirectories within the server's data directory. This abstraction allows for easier testing and potential future implementations using different storage backends.

## Usage

For usage examples, see [root.go](root.go) and [root_test.go](root_test.go).

## Interface

### Root

The `Root` interface provides methods to access subdirectories:

- `MapsDir() string` - Returns the path to the maps data directory

## Implementations

### osRoot

The default implementation (`osRoot`) uses the operating system's filesystem and is created via `NewOSRoot(baseDir string)`.

## Testing

The package includes comprehensive tests. To run them:

```bash
go test ./internal/data/...
```

## Future Extensions

The `Root` interface can be extended to include additional subdirectories as needed:

- `UsersDir()` - For user data
- `SettingsDir()` - For server settings
- etc.
