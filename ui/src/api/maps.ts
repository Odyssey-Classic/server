export interface MapLinks {
    north?: number
    east?: number
    south?: number
    west?: number
}

export interface MapTile {
    passable?: boolean
    graphics?: { graphic_id: number; properties?: Record<string, string> }[]
    blocked_directions?: { direction: number; block_inbound: boolean; block_outbound: boolean }[]
    warp?: { map_id: number; x: number; y: number } | null
    triggers?: { event: string; payload?: Record<string, string> }[]
    attributes?: Record<string, string>
}

export interface GameMap {
    id: number
    name: string
    tags: string[]
    attributes?: Record<string, string>
    last_updated: string
    version: number
    tiles: MapTile[][] | any // server sends 17x17 array; we won't rely on it yet in editor
    links: MapLinks
}

const JSON_HEADERS: HeadersInit = {
    'Content-Type': 'application/json',
}

function ok(res: Response) {
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return res
}

export async function listMaps(query?: string): Promise<GameMap[]> {
    const u = new URL('/admin/maps', window.location.origin)
    if (query) u.searchParams.set('q', query)
    const res = await fetch(u.toString(), { method: 'GET' }).then(ok)
    return res.json()
}

export async function getMap(id: number): Promise<GameMap> {
    const res = await fetch(`/admin/maps/${id}`, { method: 'GET' }).then(ok)
    return res.json()
}

export async function createMap(name: string): Promise<GameMap> {
    // minimal payload: only name; server fills defaults via NewMap
    const res = await fetch('/admin/maps', {
        method: 'POST',
        headers: JSON_HEADERS,
        body: JSON.stringify({ name }),
    }).then(ok)
    return res.json()
}

export async function updateMap(map: GameMap): Promise<void> {
    await fetch(`/admin/maps/${map.id}`, {
        method: 'PUT',
        headers: JSON_HEADERS,
        body: JSON.stringify(map),
    }).then(ok)
}
