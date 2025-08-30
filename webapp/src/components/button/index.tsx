const Button = ({ children, className }: { children: React.ReactNode, className?: string }) => {
    return (
        <button className={className}>
            {children}
        </button>
    )
}

export default Button;