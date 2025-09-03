import { useState, useRef, useEffect } from "react"
import { suggestApi } from "../../api/suggestApi"
import { PAGE_LIMIT } from "../../utils/constants"

const SearchBar = ({
  className,
  value,
  handleChange,
  resultsBoxClassName,
  resultClassName,
}: {
  className: string
  value: string
  handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void
  resultsBoxClassName?: string
  resultClassName?: string
}) => {
  const [suggestions, setSuggestions] = useState<string[]>([])
  const [searchOpen, setSearchOpen] = useState(false)
  const searchBarRef = useRef<HTMLDivElement>(null)

  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchBarRef.current && !searchBarRef.current.contains(event.target as Node)) {
        setSearchOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [])

  const onChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    handleChange(e)
    // get the redis suggestions here
    try {
      const suggestions = await suggestApi(e.target.value)
      setSuggestions(suggestions)
      setSearchOpen(suggestions.length > 0)
    } catch (error) {
      console.error(error)
      setSearchOpen(false)
    }
  }

  const handleClick = () => {
    setSearchOpen(true)
  }
  return (
    <div className="relative" ref={searchBarRef}>
    <input type="text" className={className} name="query" value={value} onChange={onChange} onClick={handleClick} />
    {searchOpen && suggestions.length > 0 && (
      <div className={`absolute top-[100%] left-0 ${suggestions.length > 3 ? 'max-h-[100px]' : ''} overflow-y-scroll suggestion overflow-x-hidden flex flex-col items-start justify-start ${resultsBoxClassName ?? ''}`}>
        {suggestions.map((suggestion) => (
          <a key={suggestion} href={`/q?query=${suggestion}&page=1&count=${PAGE_LIMIT}`} className={`py-1 px-2 hover:cursor-pointer hover:bg-gray-100 border-b border-[#676767] ${resultClassName ?? ''}`}>{suggestion}</a>
        ))}
      </div>
    )}
    </div>
  )
}

export default SearchBar