import { API_URL } from "../utils/constants"

/**
 * Suggest a query
 * @param prefix - The prefix to suggest
 * @returns 
 */
export const suggestApi = async (prefix: string) => {
    if (prefix === '' || prefix.length < 3 || prefix.includes('http') || prefix.includes('/')) {
        return []
    }
    const res = await fetch(`${API_URL}/api/suggest?prefix=${prefix}`, {
        headers: {
            'Content-Type': 'application/json',
        },
        method: 'GET',
    })
    if (!res.ok) {
        throw new Error('Error: ' + res.statusText)
    }
    return res.json()
}