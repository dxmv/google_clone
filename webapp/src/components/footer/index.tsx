import { Link, useNavigate, useLocation } from "react-router";
import { GithubIcon } from "lucide-react";

const Footer = () => {
    const location = useLocation()
    const isHomePage = location.pathname === '/'
    
    return (
        <footer className={`flex flex-row items-center justify-end gap-4 px-4 py-2 ${!isHomePage ? 'border-t-2 border-[#676767]' : ''}`}>
            
            <Link to="https://github.com/dxmv" target="_blank"><GithubIcon color="#1A54CB"/></Link>
        </footer>
    )
}

export default Footer;