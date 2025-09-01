import ResultLink from "../components/result-link";
import type { SearchResult as SearchResultType } from "../types";
import { useNavigate, useSearchParams } from "react-router";

const AllResults = ({results, currentPage, totalPages, suggestion, count}: {results: SearchResultType[], currentPage: number, totalPages: number, suggestion: string | null, count: string | null}) => {
  const navigate = useNavigate()
  const searchParams = useSearchParams()
  const query = searchParams[0].get('query')
  console.log(results)

        return (
    <>
    {/* Results */}
    <div className="px-44 py-8 flex flex-col items-start justify-start gap-4 overflow-x-hidden">
    {suggestion && <div className="italic mb-4 text-black">Did you mean:<a className="font-bold text-[#1A54CB] hover:cursor-pointer hover:underline" onClick={() => navigate(`/q?query=${suggestion}&page=1&count=${count}`)}>{suggestion}</a>?</div>}
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
  </>
  )
}

export default AllResults