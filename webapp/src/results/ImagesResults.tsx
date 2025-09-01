import { useNavigate } from "react-router";
import ImageResult from "../components/image-result";
import type { SearchResult } from "../types";

const ImagesResults = ({results, suggestion, count}: {results: SearchResult[], suggestion: string | null, count: number}) => {
    const navigate = useNavigate();

    return (
    <>
      {/* Images */}
      <div className="p-8">
        {suggestion && <div className="italic mb-4 text-black ">Did you mean: <a className="font-bold text-[#1A54CB] hover:cursor-pointer hover:underline" onClick={() => navigate(`/q?query=${suggestion}&page=1&count=${count}`)}>{suggestion}</a>?</div>}
        <div className=" grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
          {results.flatMap((result) => 
            result.doc.images
              .filter(image => image && !image.includes(".svg")) // Filter out undefined/empty strings
              .map((image, index) => (
                <ImageResult key={`${result.doc.hash}-${index}`} image={image} title={result.doc.title} url={result.doc.url} />
              ))
          )}
        </div>
      </div>
    </>
    )
}

export default ImagesResults