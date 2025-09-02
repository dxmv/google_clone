import { useSearchParams } from 'react-router'
import { useEffect, useState } from 'react'
import type { SearchResult } from '../types'
import Loading from './loading'
import Error from './error'
import SearchResults from './SearchResults'
import { searchApi } from '../api/searchApi'

interface FinalResults {
  results: SearchResult[];
  count: number;
  suggestion: string | null;
}

function index() {
  const [searchParams] = useSearchParams()
  const [results, setResults] = useState<FinalResults>({results: [
  ], count: 0, suggestion: null})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const query: string = searchParams.get('query') || ''
  const page: number = parseInt(searchParams.get('page') || '1')
  const count: number = parseInt(searchParams.get('count') || '24')

  useEffect(() => {
    const fetchResults = async () => {
        try {
          const data = await searchApi(query, page, count)
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