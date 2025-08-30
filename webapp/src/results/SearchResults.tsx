import type { SearchResult } from '../types'
import Layout from '../components/layout/layout'
import Button from '../components/button'
import { useState } from 'react'
import ResultLink from '../components/result-link'
import type { SearchResult as SearchResultType } from '../types'
import { Link, useNavigate, useSearchParams } from 'react-router'

type Tab = 'All' | 'Images'

function SearchResults({results, currentPage, totalPages}: {results: SearchResult[], currentPage: number, totalPages: number}) {
  const [tab, setTab] = useState<Tab>('All')
  const navigate = useNavigate()
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
      <div className="px-44 py-8 flex flex-col items-start justify-start gap-4 overflow-x-hidden">
        {results.map((result) => (
          <ResultLink key={result.doc.hash} result={result as SearchResultType} />
        ))}
      </div>
      {/* Pagination buttons */}
      <div className='px-44 flex flex-row items-center justify-start gap-2 text-[#1A54CB] mb-4'>
        {currentPage > 1 && (
          <a 
            className='hover:cursor-pointer hover:text-[#676767]' 
            onClick={() => navigate(`/q?query=${query}&page=${currentPage - 1}&count=${count}`)}
          >
            &lt; Previous 
          </a>
        )}
        
        {/* Generate page numbers dynamically */}
        {Array.from({ length: Math.min(totalPages, 10) }, (_, i) => {
          const pageNum = i + 1;
          const isCurrentPage = pageNum === currentPage;
          
          return (
            <a
              key={pageNum}
              className={`hover:cursor-pointer ${
                isCurrentPage 
                  ? 'font-bold underline' 
                  : 'hover:underline'
              }`}
              onClick={() => navigate(`/q?query=${query}&page=${pageNum}&count=${count}`)}
            >
              {pageNum}
            </a>
          );
        })}
        
        {currentPage < totalPages && (
          <a 
            className='hover:cursor-pointer hover:text-[#676767]' 
            onClick={() => navigate(`/q?query=${query}&page=${currentPage + 1}&count=${count}`)}
          >
            Next &gt;
          </a>
        )}
      </div>
    </Layout>
  )
}

const TabButton = ({tab, setTab, activeTab}: {tab: Tab, setTab: (tab: Tab) => void, activeTab: Tab}) => {
  return (
    <div className={`border-r-2 border-[#676767] hover:cursor-pointer hover:text-[#676767] px-4 ${tab === activeTab ? 'border-b-2 border-[#676767]' : ''}`} onClick={() => setTab(tab)}>{tab}</div>
  )
}


export default SearchResults