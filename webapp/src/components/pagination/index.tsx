import { useNavigate } from "react-router";

const VISIBLE_PAGES = 10;

const Pagination = ({
  currentPage,
  totalPages,
  query,
  count,
}: {
  currentPage: number;
  totalPages: number;
  query: string | null;
  count: string | null;
}) => {
  const navigate = useNavigate();
  
  // Calculate the range of page numbers to display
  const getPageNumbers = () => {
    const pages: number[] = [];
    const halfVisible = Math.floor(VISIBLE_PAGES / 2);
    
    let startPage = Math.max(1, currentPage - halfVisible);
    let endPage = Math.min(totalPages, currentPage + halfVisible);
    
    // Adjust if we're near the beginning or end
    if (endPage - startPage + 1 < VISIBLE_PAGES) {
      if (startPage === 1) {
        endPage = Math.min(totalPages, startPage + VISIBLE_PAGES - 1);
      } else if (endPage === totalPages) {
        startPage = Math.max(1, endPage - VISIBLE_PAGES + 1);
      }
    }
    
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }
    
    return pages;
  };
  
  const pageNumbers = getPageNumbers();
  
  return (
    <>
      {/* Pagination buttons */}
      <div className="px-44 flex flex-row items-center justify-start gap-2 text-[#1A54CB] mb-4">
        {currentPage > 1 && (
          <PreviousPageButton currentPage={currentPage} query={query} count={count} />
        )}

        {/* Page numbers */}
        {pageNumbers.map((pageNum) => (
          <a
            key={pageNum}
            className={
              pageNum === currentPage
                ? "font-bold underline"
                : "hover:cursor-pointer hover:text-[#676767]"
            }
            onClick={() => navigate(`/q?query=${query}&page=${pageNum}&count=${count}`)}
          >
            {pageNum}
          </a>
        ))}

        {currentPage < totalPages && (
          <NextPageButton currentPage={currentPage} query={query} count={count} />
        )}
      </div>
    </>
  );
};

const NextPageButton = ({currentPage, query, count}: {currentPage: number, query: string | null, count: string | null}) => {
  const navigate = useNavigate();
  return (
    <a className="hover:cursor-pointer hover:text-[#676767]" onClick={() => navigate(`/q?query=${query}&page=${currentPage + 1}&count=${count}`)}>Next &gt;</a>
  );
};

const PreviousPageButton = ({currentPage, query, count}: {currentPage: number, query: string | null, count: string | null}) => {
  const navigate = useNavigate();
  return (
    <a className="hover:cursor-pointer hover:text-[#676767]" onClick={() => navigate(`/q?query=${query}&page=${currentPage - 1}&count=${count}`)}>&lt; Previous</a>
  );
};


export default Pagination;
