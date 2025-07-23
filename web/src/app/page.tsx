'use client'

import { useState } from 'react'
import SubmitEvent from '@/components/SubmitEvent';
import ViewEvents from '@/components/ViewEvents';


type TabType = 'submit' | 'view';

export default function Home() {
  const [activeTab, setActiveTab] = useState<TabType>('submit')

  return (
    <main className="p-6 max-w-md mx-auto">
      <h1 className="text-2xl mb-4">Event Tracker</h1>

      {/* Tab Navigation */}
      <div className="flex mb-6 border-b">
        <button
          onClick={() => setActiveTab('submit')}
          className={`px-4 py-2 font-medium ${activeTab === 'submit'
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 hover:text-gray-800'
            }`}
        >
          Submit Event
        </button>
        <button
          onClick={() => setActiveTab('view')}
          className={`px-4 py-2 font-medium ${activeTab === 'view'
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 hover:text-gray-800'
            }`}
        >
          View Events
        </button>
      </div>

      {/* Content based on active tab */}
      {activeTab === 'submit' ? (
        <div>
          <SubmitEvent />
        </div>
      ) : (
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">View Events</h2>
          <ViewEvents />
        </div>
      )}
    </main>
  );
}
