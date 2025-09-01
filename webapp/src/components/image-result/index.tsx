const ImageResult = ({image, title, url}: {image: string, title: string, url: string}) => {
  return (
    <a href={url} target="_blank" rel="noopener noreferrer" className="group">
      <div className="overflow-hidden border border-gray-200">
        <img 
          src={image} 
          alt={title} 
          className="w-full h-48 object-cover group-hover:opacity-75 transition-opacity"
        />
      </div>
      <div className="mt-2">
        <h3 className="text-sm text-gray-700 truncate group-hover:underline">{title}</h3>
        <p className="text-xs text-gray-500 truncate">{url}</p>
      </div>
    </a>
  )
}

export default ImageResult      