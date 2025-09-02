import { API_URL } from "../utils/constants";

/**
 * Search for a query
 * @param query - The query to search for
 * @param page - The page number
 * @param count - The number of results to return
 * @returns 
 */
export const searchApi = async (query: string, page: number, count: number) => {
        if (query === '') {
                throw new Error('Error:No query')
                return
        }
        const res = await fetch(`${API_URL}/api/search`, {
                headers: {
                        'Content-Type': 'application/json',
                },
                method: 'POST',
                body: JSON.stringify({ query, page, count }),
        });
        
        if (!res.ok) {
                alert('Error: ' + res.statusText)
                return
        }
        
        return res.json()
}