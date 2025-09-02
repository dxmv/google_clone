import { Link, useNavigate } from "react-router"
import Button from "../../components/button"
import SearchBar from "../../components/search-bar"
import { useState, type ButtonHTMLAttributes } from "react"
import { PAGE_LIMIT } from "../../utils/constants"
const Header = ({initialQuery}: {initialQuery: string}) => {
        const [query, setQuery] = useState(initialQuery)
        const navigate = useNavigate()

        const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
                setQuery(e.target.value)
        }


        const handleSubmit = () => {
                navigate(`/q?query=${query}&page=1&count=${PAGE_LIMIT}`)
        }

        return (
                <header className="flex flex-row items-center justify-start p-8">
                        <Link to="/"><img src="/logo.png" alt="Logo" className="w-28 mr-8" /></Link>
                        <SearchBar className="w-[400px] px-2 mr-2" handleChange={handleChange} value={query} />
                        <Button className="min-w-[100px]" onClick={handleSubmit}>Search</Button>
                </header>
        )
}

export default Header