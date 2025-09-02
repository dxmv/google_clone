import { useState } from "react"
import { suggestApi } from "../../api/suggestApi"

const SearchBar = ({className, resultClassName, value, handleChange}: {className: string, resultClassName?: string, value: string, handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void}) => {
  const [suggestions, setSuggestions] = useState<string[]>([])
  const onChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    setSuggestions([])
    handleChange(e)
    // get the redis suggestions here
    try {
      const suggestions = await suggestApi(e.target.value)
      setSuggestions(suggestions)
    } catch (error) {
      console.error(error)
    }
  }
  return (
    <div className="relative">
    <input type="text" className={className} name="query" value={value} onChange={onChange} />
    {suggestions.length > 0 && (
      <div className={`absolute top-[100%] left-0 w-full bg-white`}>
        {suggestions.map((suggestion) => (
          <div  key={suggestion} className={`suggestion hover:cursor-pointer hover:bg-gray-100 border-b border-[#676767] ${resultClassName}`}>{suggestion}</div>
        ))}
      </div>
    )}
    </div>
  )
}

export default SearchBar