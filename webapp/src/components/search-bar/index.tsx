
const SearchBar = ({className, value, handleChange}: {className: string, value: string, handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void}) => {
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleChange(e)
    // get the redis suggestions here
  }
  return (
    <input type="text" className={className} name="query" value={value} onChange={onChange} />
  )
}

export default SearchBar