import { useState } from "react"

// Define the form data type
type FormData = {
  query: string;
}

const API_URL = import.meta.env.VITE_QUERY_API_URL || 'https://localhost:8000'

function App() {
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
    const res = await fetch(`${API_URL}/api/search`, {
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
      body: JSON.stringify({ ...formData }),
    });

    if (!res.ok) {
      alert('Error: ' + res.statusText)
      return
    }

    const data = await res.json()
    alert(data.message)
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
