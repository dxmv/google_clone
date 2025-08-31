import { useSearchParams } from 'react-router'
import { useEffect, useState } from 'react'
import type { SearchResult } from '../types'
import Loading from './Loading'
import Error from './error'
import SearchResults from './SearchResults'
        
const API_URL = import.meta.env.VITE_QUERY_API_URL || 'https://localhost:8000'

interface FinalResults {
  results: SearchResult[];
  count: number;
  suggestion: string | null;
}

function index() {
  const [searchParams] = useSearchParams()
  const [results, setResults] = useState<FinalResults>({results: [
  ], count: 0, suggestion: null})
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const query: string = searchParams.get('query') || ''
  const page: number = parseInt(searchParams.get('page') || '1')
  const count: number = parseInt(searchParams.get('count') || '24')

  useEffect(() => {
    const fetchResults = async () => {
        try {
                if (query === '') {
                        setError('No query')
                        setLoading(false)
                        return
                }
                const res = await fetch(`${API_URL}/api/search`, {
                        headers: {
                                'Content-Type': 'application/json',
                        },
                        method: 'POST',
                        body: JSON.stringify({ query, page, count }),
                });
                
                if (!res.ok) {
                  alert('Error: ' + res.statusText)
                  return
                }
                
                const data = await res.json()
                console.log(data)
                setResults({results: data.results, count: data.total, suggestion: data.suggestion})
                setLoading(false)
                setError(null)
    } catch (err) {
        setError(err && typeof err === 'object' && 'message' in err ? (err as Error).message : 'An unknown error occurred')
        setLoading(false)
    }
    }
    fetchResults()
  }, [query, page, count])

  return (
    <>
      {loading ? (
        <Loading />
      ) : error || query === '' ? (
        <Error error={error}/>
      ) : (
        <>
          <SearchResults results={results.results} currentPage={page} totalPages={Math.ceil(results.count / count)} suggestion={results.suggestion} />
        </>
      )}
    </>
  )
}

export default index