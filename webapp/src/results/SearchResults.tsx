import type { SearchResult } from '../types'
import Layout from '../components/layout/layout'
import Button from '../components/button'
import { useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import AllResults from './AllResults'
import ImagesResults from './ImagesResults'

type Tab = 'All' | 'Images'

function SearchResults({results, currentPage, totalPages, suggestion}: {results: SearchResult[], currentPage: number, totalPages: number, suggestion: string | null}) {
  const [tab, setTab] = useState<Tab>('All')
  const searchParams = useSearchParams()
  const query = searchParams[0].get('query')
  const count = searchParams[0].get('count')
  return (
    <Layout>
      {/* Header */}
      <header className="flex flex-row items-center justify-start p-8">
        <Link to="/"><img src="/logo.png" alt="Logo" className="w-28 mr-8" /></Link>
        <input type="text" className="w-[400px] px-2 mr-2" />
        <Button className="min-w-[100px]">Search</Button>
      </header>
      {/* Tabs */}
      <div className="flex flex-row items-center justify-start border-b-2 border-[#676767] px-44">
        <TabButton tab="All" setTab={setTab} activeTab={tab} />
        <TabButton tab="Images" setTab={setTab} activeTab={tab} />
      </div>
      {/* Results */}
      {tab === 'All' ? (
        <AllResults results={results} currentPage={currentPage} totalPages={totalPages} suggestion={suggestion} count={count} />
      ) : <ImagesResults results={results}  suggestion={suggestion} />}
    </Layout>
  )
}

const TabButton = ({tab, setTab, activeTab}: {tab: Tab, setTab: (tab: Tab) => void, activeTab: Tab}) => {
  return (
    <div className={`border-r-2 border-[#676767] hover:cursor-pointer hover:text-[#676767] px-4 ${tab === activeTab ? 'border-b-2 border-[#676767]' : ''}`} onClick={() => setTab(tab)}>{tab}</div>
  )
}


export default SearchResults