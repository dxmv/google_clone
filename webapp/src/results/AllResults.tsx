import Pagination from "../components/pagination";
import ResultLink from "../components/result-link";
import type { SearchResult as SearchResultType } from "../types";
import { useNavigate, useSearchParams } from "react-router";

const AllResults = ({
  results,
  currentPage,
  totalPages,
  suggestion,
  count,
  query_time,
}: {
  results: SearchResultType[];
  currentPage: number;
  totalPages: number;
  suggestion: string | null;
  count: string | null;
  query_time: number;
}) => {
  const navigate = useNavigate();
  const searchParams = useSearchParams();
  const query = searchParams[0].get("query");
  console.log(results);

  return (
    <>
      {/* Results */}
      <div className="px-44 py-8 flex flex-col items-start justify-start gap-4 overflow-x-hidden">
      {!suggestion && (<p className="italic text-[#676767]">Query time: {query_time} seconds</p>)}
      {suggestion && (
          <div className="italic text-black">
            Did you mean:
            <a
              className="font-bold text-[#1A54CB] hover:cursor-pointer hover:underline"
              onClick={() =>
                navigate(`/q?query=${suggestion}&page=1&count=${count}`)
              }
            >
              {suggestion}
            </a>
            ?
          </div>
        )}
        {results.map((result) => (
          <ResultLink
            key={result.doc.hash}
            result={result as SearchResultType}
          />
        ))}
      </div>
      <Pagination
        currentPage={currentPage}
        totalPages={totalPages}
        query={query}
        count={count}
      />
    </>
  );
};

export default AllResults;
