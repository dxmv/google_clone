import type { SearchResult } from '../types'

function SearchResults({results}: {results: SearchResult[]}) {
  return (
    <div>
        {results.map((result) => (
            <a key={result.doc.hash} href={result.doc.url}>
                <h2>{result.doc.title}</h2>
            </a>
        ))}
    </div>
  )
}

export default SearchResults