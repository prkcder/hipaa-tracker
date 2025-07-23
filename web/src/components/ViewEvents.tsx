'use client'

import { useEffect, useState } from 'react'

type Payload = Record<string, string | number | boolean | null | undefined>;

type Event = {
    id: number
    event_type: string
    sanitized: boolean
    created_at: string
    payload: Payload
}

export default function ViewEvents() {
    const [events, setEvents] = useState<Event[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [query, setQuery] = useState('')
    const [showRedacted, setShowRedacted] = useState(true)
    const [autoRefresh, setAutoRefresh] = useState(false)
    const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null)
    const [flags, setFlags] = useState<Record<number, boolean>>({})

    // Load flags from localStorage
    useEffect(() => {
        const saved = localStorage.getItem('flaggedEvents')
        if (saved) {
            setFlags(JSON.parse(saved))
        }
    }, [])

    // Persist flags when they change
    useEffect(() => {
        localStorage.setItem('flaggedEvents', JSON.stringify(flags))
    }, [flags])

    const fetchEvents = async () => {
        try {
            const res = await fetch(process.env.NEXT_PUBLIC_BACKEND_URL + '/events')
            if (!res.ok) throw new Error(`HTTP ${res.status}`)
            const data = await res.json()
            setEvents(data)
        } catch (err) {
            setError('Failed to load events')
            console.error(err)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchEvents()
        if (autoRefresh) {
            const interval = setInterval(fetchEvents, 5000)
            setRefreshInterval(interval)
            return () => clearInterval(interval)
        }
        return () => {
            if (refreshInterval) clearInterval(refreshInterval)
        }
    }, [autoRefresh])

    const filtered = events.filter((e) =>
        JSON.stringify(e).toLowerCase().includes(query.toLowerCase())
    )

    const maybeRedact = (payload: Payload) => {
        if (showRedacted) return payload
        const redacted = { ...payload }
        for (const key of ['email', 'phone', 'dob', 'mrn', 'insurance_id', 'device_id', 'zipcode']) {
            if (redacted[key]) redacted[key] = '🔒 REDACTED'
        }
        return redacted
    }

    const copyJSON = (e: Event) => {
        navigator.clipboard.writeText(JSON.stringify(e, null, 2))
        alert('Copied JSON to clipboard!')
    }

    const exportCSV = () => {
        const headers = ['id', 'event_type', 'sanitized', 'created_at', 'payload']
        const rows = filtered.map((e) => [
            e.id,
            e.event_type,
            e.sanitized,
            new Date(e.created_at).toISOString(),
            JSON.stringify(maybeRedact(e.payload)).replace(/"/g, '""'),
        ])
        const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n')
        const blob = new Blob([csv], { type: 'text/csv' })
        const url = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = 'events.csv'
        a.click()
        URL.revokeObjectURL(url)
    }

    const toggleFlag = (id: number) => {
        setFlags((prev) => ({ ...prev, [id]: !prev[id] }))
    }

    return (
        <div className="space-y-6 mt-6">
            <div className="flex flex-col sm:flex-row gap-4 items-center">
                <input
                    className="border px-3 py-1 rounded w-full sm:w-1/2"
                    placeholder="Search by name, email, etc..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                />
                <label className="flex items-center gap-2 text-sm">
                    <input type="checkbox" checked={autoRefresh} onChange={() => setAutoRefresh(!autoRefresh)} />
                    Auto-refresh
                </label>
                <label className="flex items-center gap-2 text-sm">
                    <input type="checkbox" checked={showRedacted} onChange={() => setShowRedacted(!showRedacted)} />
                    Show raw
                </label>
                <button onClick={exportCSV} className="text-sm bg-green-600 text-white px-3 py-1 rounded">
                    📤 Export CSV
                </button>
            </div>

            {loading ? (
                <p>Loading events...</p>
            ) : error ? (
                <p className="text-red-500">{error}</p>
            ) : filtered.length === 0 ? (
                <p>No matching events found.</p>
            ) : (
                filtered.map((event) => (
                    <div key={event.id} className="p-4 border rounded shadow-sm bg-white relative">
                        <div className="absolute top-2 right-2 flex gap-2">
                            <button
                                onClick={() => copyJSON(event)}
                                className="text-xs text-blue-600 hover:underline"
                            >
                                📋 Copy JSON
                            </button>
                            <button
                                onClick={() => toggleFlag(event.id)}
                                className={`text-xs ${flags[event.id] ? 'text-red-600' : 'text-gray-400'}`}
                            >
                                🚩 {flags[event.id] ? 'Flagged' : 'Flag'}
                            </button>
                        </div>
                        <p><strong>ID:</strong> {event.id}</p>
                        <p><strong>Type:</strong> {event.event_type}</p>
                        <p><strong>Sanitized:</strong> {event.sanitized ? 'Yes' : 'No'}</p>
                        <p><strong>Timestamp:</strong> {new Date(event.created_at).toLocaleString()}</p>
                        <details className="mt-2">
                            <summary className="cursor-pointer text-blue-600">Payload</summary>
                            <pre className="bg-gray-100 p-2 rounded mt-2 text-sm overflow-x-auto">
                                {JSON.stringify(maybeRedact(event.payload), null, 2)}
                            </pre>
                        </details>
                    </div>
                ))
            )}
        </div>
    )
}