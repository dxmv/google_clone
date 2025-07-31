/// <reference types="vite/client" />

export interface SearchResult {
  doc: Doc;
  score: number;
  term_count: number;
}

export interface Doc{
        url: string;
        depth: number;
        title: string;
        hash: string;
        links: string[];
}
