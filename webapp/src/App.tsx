import { useState } from "react"
import { useNavigate } from "react-router"

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
    <div className="h-screen w-screen bg-white flex flex-col items-center justify-center">
      <h1 className="text-4xl font-bold">Hello World</h1>
      <form onSubmit={handleSubmit} className="flex flex-col items-center justify-center gap-4 mt-8">
          <input type="text" className="border-2 border-gray-300 rounded-md p-2" name="query" value={formData.query} onChange={handleChange} />
          <button type="submit" className="bg-blue-500 text-white rounded-md p-2">Submit</button>
      </form>

    </div>
  )
}

export default App
