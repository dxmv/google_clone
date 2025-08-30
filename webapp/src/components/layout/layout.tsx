import Footer from "../footer";

const Layout = ({ children }: { children: React.ReactNode }) => {
    return (
        <div className="flex flex-col min-h-screen w-full overflow-x-hidden"
        style={{
          backgroundImage: 'url(/bg.png)',
          backgroundRepeat: 'repeat',
        }}>
          {children}
          <Footer />
        </div>
    )
}

export default Layout;