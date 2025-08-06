'use client';

import { useEffect, useState } from 'react';

type Payload = Record<string, string | number | boolean | null | undefined>;

type Event = {
    id: number
    event_type: string
    sanitized: boolean
    created_at: string
    payload: Payload
};

export default function ViewEvents() {
    const [events, setEvents] = useState<Event[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [query, setQuery] = useState('');
    // const [showRedacted, setShowRedacted] = useState(true);
    const [autoRefresh, setAutoRefresh] = useState(false);
    const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null);
    const [flags, setFlags] = useState<Record<number, boolean>>({});
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 5;

    const filter = events.filter((e) =>
        JSON.stringify(e).toLowerCase().includes(query.toLowerCase())
    );

    const totalPages = Math.ceil(filter.length / itemsPerPage);

    const paginated = filter.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage
    )


    // Load flags from localStorage
    useEffect(() => {
        const saved = localStorage.getItem('flaggedEvents')
        if (saved) {
            setFlags(JSON.parse(saved))
        }
    }, []);

    // Persist flags when they change
    useEffect(() => {
        localStorage.setItem('flaggedEvents', JSON.stringify(flags))
    }, [flags]);

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
    };

    useEffect(() => {
        fetchEvents();
        if (autoRefresh) {
            const interval = setInterval(fetchEvents, 5000);
            setRefreshInterval(interval);
            return () => clearInterval(interval);
        }
        return () => {
            if (refreshInterval) clearInterval(refreshInterval);
        }
    }, [autoRefresh, refreshInterval]);

    const maybeRedact = (payload: Payload) => {
        // if (showRedacted) return payload;
        const redacted = { ...payload };
        // for (const key of ['email', 'phone', 'dob', 'mrn', 'insurance_id', 'device_id', 'zipcode']) {
        //     if (redacted[key]) redacted[key] = '🔒 REDACTED'
        // };
        return redacted;
    };

    const copyJSON = (e: Event) => {
        navigator.clipboard.writeText(JSON.stringify(e, null, 2));
        alert('Copied JSON to clipboard!');
    };

    const exportCSV = () => {
        const headers = ['id', 'event_type', 'sanitized', 'created_at', 'payload'];
        const rows = filter.map((e) => [
            e.id,
            e.event_type,
            e.sanitized,
            new Date(e.created_at).toISOString(),
            JSON.stringify(maybeRedact(e.payload)).replace(/"/g, '""'),
        ]);
        const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n');
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'events.csv';
        a.click();
        URL.revokeObjectURL(url);
    };

    const toggleFlag = (id: number) => {
        setFlags((prev) => ({ ...prev, [id]: !prev[id] }));
    };

    return (
        <div className="space-y-6 mt-6">
            <div className="flex flex-col sm:flex-row gap-4 items-center">
                <input
                    className="border dark:border-gray-600 px-3 py-1 rounded w-full sm:w-1/2 
              text-gray-900 dark:text-white 
              bg-white dark:bg-gray-800 
              placeholder-gray-400 dark:placeholder-gray-500"
                    placeholder="Search by name, email, etc..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                />
                <label className="flex items-center gap-2 text-sm">
                    <input type="checkbox" checked={autoRefresh} onChange={() => setAutoRefresh(!autoRefresh)} />
                    Auto-refresh
                </label>
                {/* <label className="flex items-center gap-2 text-sm">
                    <input type="checkbox" checked={showRedacted} onChange={() => setShowRedacted(!showRedacted)} />
                    Show raw
                </label> */}
                <button onClick={exportCSV} className="text-sm bg-green-600 text-white px-3 py-1 rounded">
                    📤 Export CSV
                </button>
            </div>

            {loading ? (
                <p>Loading events...</p>
            ) : error ? (
                <p className="text-red-500">{error}</p>
            ) : filter.length === 0 ? (
                <p>No matching events found.</p>
            ) : (
                <>
                    {paginated.map((event) => (
                        <div key={event.id} className="p-4 border dark:border-gray-600 rounded shadow-sm bg-gray-50 dark:bg-gray-800 relative">
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
                                <pre className="bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 p-3 rounded mt-2 text-sm overflow-x-auto border dark:border-gray-600">
                                    {JSON.stringify(maybeRedact(event.payload), null, 2)}
                                </pre>
                            </details>
                        </div>
                    ))}

                    {/* Pagination Controls */}
                    <div className="flex justify-center items-center mt-6 space-x-2 text-sm">
                        <button
                            onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
                            className="px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 disabled:opacity-50"
                            disabled={currentPage === 1}
                        >
                            ◀
                        </button>

                        {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                            <button
                                key={page}
                                onClick={() => setCurrentPage(page)}
                                className={`px-3 py-1 rounded ${currentPage === page
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-gray-100 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600'
                                    }`}
                            >
                                {page}
                            </button>
                        ))}

                        <button
                            onClick={() => setCurrentPage((p) => Math.min(p + 1, totalPages))}
                            className="px-2 py-1 rounded bg-gray-200 dark:bg-gray-700 disabled:opacity-50"
                            disabled={currentPage === totalPages}
                        >
                            ▶
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}