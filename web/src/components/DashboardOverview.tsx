'use client'

import { useEffect, useState } from 'react';

type Event = {
    id: number,
    event_type: string,
    sanitized: boolean,
    created_at: string
};


export default function Dashboard() {

    const [events, setEvents] = useState<Event[]>([]);
    const [flaggedCount, setFlaggedCount] = useState(0);

    useEffect(() => {
        fetch(process.env.NEXT_PUBLIC_BACKEND_URL + '/events')
            .then((res) => res.json())
            .then((data) => setEvents(data))
            .catch(console.error)
    }, []);

    useEffect(() => {
        // Safe: only runs in the browser
        const flagged = JSON.parse(localStorage.getItem('flaggedEvents') || '{}');
        const count = Object.values(flagged).filter(Boolean).length;
        setFlaggedCount(count);
    }, []);

    const total = events.length;
    const sanitized = events.filter((e) => e.sanitized).length;


    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold">Dashboard Overview</h2>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div className="bg-blue-100 dark:bg-blue-900 p-4 rounded">
                    <p className="text-sm">Total Events</p>
                    <p className="text-xl font-bold">{total}</p>
                </div>
                <div className="bg-green-100 dark:bg-green-900 p-4 rounded">
                    <p className="text-sm">Sanitized</p>
                    <p className="text-xl font-bold">{sanitized}</p>
                </div>
                <div className="bg-red-100 dark:bg-red-900 p-4 rounded">
                    <p className="text-sm">Flagged</p>
                    <p className="text-xl font-bold">{flaggedCount}</p>
                </div>
            </div>

            {/* Recent Events Table */}
            <div className="bg-white dark:bg-gray-800 p-4 rounded shadow">
                <h3 className="text-lg font-semibold mb-2">Recent Events</h3>
                <table className="w-full text-sm table-fixed">
                    <thead>
                        <tr className="bg-gray-100 dark:bg-gray-700">
                            <th className="text-left px-4 py-2 w-1/12">ID</th>
                            <th className="text-left px-4 py-2 w-1/12">Type</th>
                            <th className="text-left px-4 py-2 w-1/12">Sanitized</th>
                            <th className="text-left px-4 py-2 w-1/12">Time</th>
                        </tr>
                    </thead>
                    <tbody>
                        {events.slice(0, 5).map((e) => (
                            <tr key={e.id}  className="border-t border-gray-200 dark:border-gray-700">
                                <td className="px-4 py-2">{e.id}</td>
                                <td className="px-4 py-2">{e.event_type}</td>
                                <td className="px-4 py-2">{e.sanitized ? 'Yes' : 'No'}</td>
                                <td className="px-4 py-2">{new Date(e.created_at).toLocaleString()}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}