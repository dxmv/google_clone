import type { SearchResult } from "../types";
import Layout from "../components/layout/layout";
import { useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import AllResults from "./AllResults";
import ImagesResults from "./ImagesResults";
import Header from "./header";

type Tab = "All" | "Images";

function SearchResults({
  results,
  currentPage,
  totalPages,
  suggestion,
  query_time,
}: {
  results: SearchResult[];
  currentPage: number;
  totalPages: number;
  suggestion: string | null;
  query_time: number;
}) {
  const [tab, setTab] = useState<Tab>("All");
  const searchParams = useSearchParams();
  const query = searchParams[0].get("query");
  const count = searchParams[0].get("count");
  const navigate = useNavigate();
  return (
    <Layout>
      <Header initialQuery={query || ""} />
      {/* Tabs */}
      <div className="flex flex-row items-center justify-start border-b-2 border-[#676767] px-44">
        <TabButton tab="All" setTab={setTab} activeTab={tab} />
        <TabButton tab="Images" setTab={setTab} activeTab={tab} />
      </div>

      {/* Results */}
      {tab === "All" ? (
        <AllResults
          results={results}
          currentPage={currentPage}
          totalPages={totalPages}
          suggestion={suggestion}
          count={count}
          query_time={query_time}
        />
      ) : (
        <ImagesResults count={0} results={results} suggestion={suggestion} />
      )}
    </Layout>
  );
}

const TabButton = ({
  tab,
  setTab,
  activeTab,
}: {
  tab: Tab;
  setTab: (tab: Tab) => void;
  activeTab: Tab;
}) => {
  return (
    <div
      className={`border-r-2 border-[#676767] hover:cursor-pointer hover:text-[#676767] px-4 ${
        tab === activeTab ? "border-b-2 border-[#676767]" : ""
      }`}
      onClick={() => setTab(tab)}
    >
      {tab}
    </div>
  );
};

export default SearchResults;
