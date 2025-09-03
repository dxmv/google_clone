export interface SearchResult {
  doc: {
    hash: string;
    url: string;
    title: string;
    score: number;
    images: string[];
    first_paragraph: string;
  };
}

export interface FinalResult {
  results: SearchResult[];
  total: number;
  suggestion: string | null;
  query_time: number;
}

