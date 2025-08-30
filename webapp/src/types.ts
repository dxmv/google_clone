export interface SearchResult {
  doc: {
    hash: string;
    url: string;
    title: string;
    score: number;
    images: string[];
  };
}
