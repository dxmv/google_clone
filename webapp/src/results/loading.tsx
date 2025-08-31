import Layout from '../components/layout/layout'

function Loading() {
  return (
    <Layout>
      <main 
      className="flex flex-col items-center justify-center flex-1"
      >
        <h1 className="text-lg mb-4">
          Loading
          <span className="loading-dots">
            <span>.</span>
            <span>.</span>
            <span>.</span>
          </span>
        </h1>
        <div className="indeterminate-loader"></div>
      </main>
    </Layout>
  )
}

export default Loading