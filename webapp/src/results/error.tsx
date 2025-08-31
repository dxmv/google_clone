
import { X } from 'lucide-react'
import Layout from '../components/layout/layout'
import { Link } from 'react-router'

function Error({error}: {error: string | null}) {
  return (
    <div className="flex flex-col items-center justify-center flex-1 h-screen w-screen overflow-hidden">
      <main 
      className="flex flex-col items-center justify-center flex-1"
      >
        <div className="flex flex-row items-center justify-center gap-4">
          <div className="flex flex-row items-center justify-center bg-[#CB261A] p-1 border-1 border-black rounded-full">
            <X className="text-white font-bold" size={24}/>
          </div>
          <div className="flex flex-col items-start justify-start">
            <h1 className="text-2xl font-bold">System Error</h1>
            <p className="text-sm">{error || 'Something went wrong'}</p>
          </div>
        </div>
        <div className="flex flex-col items-center justify-center mt-8 gap-4">
          <Link to="/" className="text-sm text-[#1A54CB] hover:cursor-pointer hover:text-[#676767]">Go back to home &gt;</Link>
          <Link className="text-sm text-[#1A54CB] hover:cursor-pointer hover:text-[#676767]" onClick={() => window.location.reload()}>Reload &gt;</Link>
        </div>
      </main>
    </div>
  )
}

export default Error