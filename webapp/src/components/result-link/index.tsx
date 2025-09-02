import type { SearchResult } from '../../types'

const MAX_FIRST_PARAGRAPH_LENGTH = 20

const ResultLink = ({result}: {result: SearchResult}) => {
  return (
    <div className="flex flex-col items-start justify-start gap-1 max-w-1/2">
      {/* Favicon and Title */}
      <div className="flex flex-row items-center justify-start gap-4">
        <img src={"https://upload.wikimedia.org/wikipedia/commons/thumb/5/5a/Wikipedia%27s_W.svg/2000px-Wikipedia%27s_W.svg.png"} alt="Favicon" 
        className="w-6 h-6" />
        <div className="flex flex-col items-start justify-start">
          <h3 className="text-sm font-bold">Wikipedia</h3>
          <p className="text-xs text-[#676767]">{result.doc.url}</p>
        </div>
      </div>
      <div className="flex flex-col items-start justify-start gap-1">
        <h2 className="text-xl text-[#1A54CB] font-bold hover:cursor-pointer hover:underline">{result.doc.title}</h2>
        <p className="text-md text-[#676767]">{result.doc.first_paragraph.split(" ").length > MAX_FIRST_PARAGRAPH_LENGTH ? result.doc.first_paragraph.split(" ").slice(0, MAX_FIRST_PARAGRAPH_LENGTH).join(" ") + '...' : result.doc.first_paragraph}</p>
      </div>
    </div>
  )
}

export default ResultLink