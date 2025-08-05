'use client';

import { useState } from 'react';

type TrackerResult = {
    tracker_name: string
    tracker_domain: string
    risk_level: string
    page_url: string
    created_at: string
    scanned_url: string
}

export default function TrackerScanner() {
    const [url, setUrl] = useState('')
    const [results, setResults] = useState<TrackerResult[] | null>(null)
    const [loading, setLoading] = useState(false)

    const normalizeUrl = (input: string): string => {
        // Trim and lowercase for consistency
        let url = input.trim();

        // If protocol is missing, assume https
        if (!/^https?:\/\//i.test(url)) {
            url = 'https://' + url
        };

        return url;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!url.trim()) return;

        setLoading(true);

        const normalizedUrl = normalizeUrl(url);
        const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8000';

        console.log('debug here: Making request to:', `${BACKEND_URL}/api/scan`);
        console.log('debug here: Request body:', { url: normalizedUrl });

        try {
            const res = await fetch(`${BACKEND_URL}/api/scan`, {
                method: 'POST',
                body: JSON.stringify({ url: normalizedUrl }),
                headers: { 'Content-Type': 'application/json' },
            });

            console.log('debug here: Response status:', res.status);
            console.log('debug here: Response ok:', res.ok);

            if (!res.ok) {
                const errorText = await res.text();
                console.log('debug here: Error response:', errorText);
                throw new Error(`Scan failed: ${errorText}`);
            }

            const data = await res.json();
            console.log('debug here: Response data:', data);
            console.log('debug here: Data type:', typeof data);
            console.log('debug here: Data is array:', Array.isArray(data));

            if (data === null) {
                console.log('debug here: Data is null - something wrong with backend');
                setResults([]);
                return;
            }

            if (!Array.isArray(data)) {
                console.log('debug here: Data is not an array:', data);
                setResults([]);
                return;
            }

            console.log('debug here: Data length:', data.length);
            setResults(data);
        } catch (error) {
            console.error('debug here: Catch error:', error);
            setResults([]);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-3xl mx-auto mt-10 px-4">
            <h1 className="text-2xl font-bold mb-4">Website Tracker Scanner</h1>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4 mb-6">
                <input
                    type="text"
                    placeholder="Enter website URL"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    className="border dark:border-gray-600 px-3 py-1 rounded w-full sm:w-1/2 
              text-gray-900 dark:text-white 
              bg-white dark:bg-gray-800 
              placeholder-gray-400 dark:placeholder-gray-500"
                />
                <button
                    type="submit"
                    className="bg-blue-600 text-white px-4 py-2 rounded"
                    disabled={loading}
                >
                    {loading ? 'Scanning...' : 'Scan Website'}
                </button>
            </form>

            {results && results.length > 0 && (
                <table className="w-full text-sm border mb-4">
                    <thead>
                        <tr>
                            <th className="border p-2">Tracker Name</th>
                            <th className="border p-2">Domain</th>
                            <th className="border p-2">Risk Level</th>
                            <th className="border p-2">Page URL</th>
                        </tr>
                    </thead>
                    <tbody>
                        {results.map((tracker, idx) => (
                            <tr key={idx}>
                                <td className="border p-2">{tracker.tracker_name}</td>
                                <td className="border p-2">{tracker.tracker_domain}</td>
                                <td className="border p-2">{tracker.risk_level}</td>
                                <td className="border p-2">{tracker.page_url}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    )
}