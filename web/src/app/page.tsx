'use client'

import { useState } from 'react'
import SubmitEvent from '@/components/SubmitEvent';
import ViewEvents from '@/components/ViewEvents';
import DashboardOverview from '@/components/DashboardOverview';
import TrackerScanner from '@/components/TrackerScanner';


type TabType = 'dashboard' | 'submit' | 'view' | 'tracker';

export default function Home() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard');


  const renderTab = () => {
    switch (activeTab) {
      case 'dashboard':
        return <DashboardOverview />
      case 'submit':
        return <SubmitEvent />
      case 'view':
        return <ViewEvents />
      case 'tracker':
        return <TrackerScanner />
    }
  };



  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* Sidebar */}
      <aside className="w-full md:w-64 bg-gray-100 dark:bg-gray-900 p-4">
        <h1 className="text-xl font-bold mb-4">HIPAA Tracker</h1>
        <nav className="space-y-2">
          <button onClick={() => setActiveTab('dashboard')} className={`block w-full text-left px-3 py-2 rounded ${activeTab === 'dashboard' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>
            📊 Dashboard
          </button>
          <button onClick={() => setActiveTab('submit')} className={`block w-full text-left px-3 py-2 rounded ${activeTab === 'submit' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>
            📝 Submit Event
          </button>
          <button onClick={() => setActiveTab('view')} className={`block w-full text-left px-3 py-2 rounded ${activeTab === 'view' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>
            📁 View Events
          </button>
          <button onClick={() => setActiveTab('tracker')} className={`block w-full text-left px-3 py-2 rounded ${activeTab === 'tracker' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>
            🌐 Web Tracker
          </button>
        </nav>
      </aside>

      {/* Main content */}
      <main className="flex-1 p-6 bg-white dark:bg-black">
        {renderTab()}
      </main>
    </div>

  )
}
