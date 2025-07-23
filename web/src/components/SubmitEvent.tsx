'use client'

import { useState } from 'react'

const defaultFormState = {
    name: '',
    email: '',
    phone: '',
    dob: '',
    zipcode: '',
    activity_type: '',
    mrn: '',
    insurance_id: '',
    device_id: '',
    timestamp: new Date().toISOString(),
}

export default function SubmitEvent() {
    const [form, setForm] = useState(defaultFormState);
    const [loading, setLoading] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!form.name || !form.email || !form.activity_type) {
            alert('Please fill out Name, Email, and Activity Type.')
            return
        }

        const data = {
            event_type: form.activity_type || 'activity', // fallback default
            payload: {
                // name: form.name,
                // email: form.email,
                ...form,
                timestamp: new Date().toISOString(),
            }
        };

        const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8000';

        try {
            const res = await fetch(`${BACKEND_URL}/event`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data),
            });

            if (!res.ok) {
                const text = await res.text();
                console.error('Server responded with error:', res.status, text);
                alert('Server error: ' + text);
                return;
            }

            alert('Event sent!');
            setForm(defaultFormState);
        } catch (err) {
            console.error('Fetch error:', err);
            if (err instanceof Error) {
                alert('Fetch error: ' + err.message);
            } else {
                alert('An unknown error occurred');
            }
        } finally {
            setLoading(false)
        }
    };

    return (
        <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-900 p-6 rounded-lg shadow space-y-4">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Submit Event</h2>
            {Object.entries(defaultFormState).map(([key, _]) => (
                key !== 'timestamp' && (
                    <input
                        key={key}
                        name={key}
                        type={key === 'dob' ? 'date' : key === 'email' ? 'email' : 'text'}
                        placeholder={key.replace('_', ' ').toUpperCase()}
                        value={form[key as keyof typeof form]}
                        onChange={handleChange}
                        className="w-full border dark:border-gray-600 rounded px-3 py-2 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white"
                    />
                )
            ))}
            <button
                type="submit"
                disabled={loading}
                className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded"
            >
                {loading ? 'Submitting...' : 'Track Event'}
            </button>
        </form>
    )
}