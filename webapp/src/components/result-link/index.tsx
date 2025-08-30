import type { SearchResult } from '../../types'

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
        <p className="text-md text-[#676767]">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
      </div>
    </div>
  )
}

export default ResultLink