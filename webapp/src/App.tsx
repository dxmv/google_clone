import { useState } from "react"
import { useNavigate } from "react-router"
import Button from "./components/button";
import { GithubIcon } from "lucide-react";

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
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    navigate(`/q?query=${formData.query}&page=1&count=24`)
  }

  return (
    <div className="flex flex-col h-screen w-screen"
    style={{
      backgroundImage: 'url(/bg.png)',
      backgroundRepeat: 'repeat',
      overflow: 'none',
    }}>
      <main 
      className="flex flex-col items-center justify-center flex-1"
      >
        <img src="/logo.png" alt="Logo" className="w-[300px]" />
        <form onSubmit={handleSubmit} className="flex flex-col items-center justify-center gap-4 mt-4">
            <input type="text" className="min-w-[600px] py-1 px-2" name="query" value={formData.query} onChange={handleChange} />
            <div className="flex flex-row items-center justify-center gap-4 mt-4">
              <Button className="min-w-[150px] py-1">Search</Button>
              <Button className="min-w-[150px] py-1">I'm Feeling Lucky</Button>
            </div>
        </form>

      </main>
      <footer className="flex flex-row items-center justify-end gap-4 px-4 pb-2">
        <GithubIcon color="#1A54CB"/>
      </footer>
    </div>
  )
}

export default App
