import { useState } from "react"
import { useNavigate } from "react-router"
import Button from "./components/button";
import Layout from "./components/layout/layout";
import SearchBar from "./components/search-bar";
import { PAGE_LIMIT } from "./utils/constants";
import { searchApi } from "./api/searchApi";

// Define the form data type
type FormData = {
  query: string;
}



function App() {
  const navigate = useNavigate()
  const [formData, setFormData] = useState<FormData>({
    query: '',
  });

  // Handle input change
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  // Handle form submission
  const handleSubmit = async () => {
    if(formData.query === '') {
      return;
    }
    navigate(`/q?query=${formData.query}&page=1&count=${PAGE_LIMIT}`)
  }

  const handleFeelingLucky = async () => {
    if(formData.query === '') {
      return;
    }
    const data = await searchApi(formData.query, 1, PAGE_LIMIT);
    if(data.results.length < 0) {
      return;
    }
    const firstResult = data.results[0];
    window.location.href = firstResult.doc.url;
  }

  return (
    <Layout>
      <main 
      className="flex flex-col items-center justify-center flex-1"
      >
        <img src="/logo.png" alt="Logo" className="w-[300px]" />
        <form onSubmit={(e)=>{e.preventDefault()}} className="flex flex-col items-center justify-center gap-4 mt-4">
            <SearchBar className="min-w-[600px] py-1 px-2" resultsBoxClassName="w-[600px]" resultClassName="w-[600px]" handleChange={handleChange} value={formData.query} />
            <div className="flex flex-row items-center justify-center gap-4 mt-4">
              <Button className="min-w-[150px] py-1" onClick={handleSubmit}>Search</Button>
              <Button className="min-w-[150px] py-1" onClick={handleFeelingLucky}>I'm Feeling Lucky</Button>
            </div>
        </form>
      </main>
    </Layout>
  )
}

export default App
