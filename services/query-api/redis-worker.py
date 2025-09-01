import redis

r = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

QUEUE_NAME = "query_queue"
SORTED_SET_NAME = "searched_ngrams"
LEX_SET_NAME = "lex_set"

def enqueue_query(query: str):
        '''
        Enqueue a query to the Redis queue
        '''
        r.rpush(QUEUE_NAME, query)

def dequeue_query()->str:
        '''
        Dequeue a query from the Redis queue
        '''
        return r.lpop(QUEUE_NAME)

def generate_ngrams(tokens: list[str]):
        '''
        Generate n-grams from a query for n=2,3
        '''
        ngrams = []
        for n in [2,3]:
                for i in range(len(tokens) - n + 1):
                        ngram = " ".join(tokens[i:i+n])
                        ngrams.append(ngram)
        return ngrams

def tokenize(query: str) -> list[str]:
    '''
    Tokenize a query
    '''
    return query.lower().split()

def process_dequeued_query(query: str):
        '''
        Process a query
        '''
        print(f"Processing query: \"{query}\"\n")
        tokenized_query = tokenize(query)
        ngrams = generate_ngrams(tokenized_query)
        print(f"Done generating n-grams\n")
        # fill redis sorted set with ngrams
        for ngram in ngrams:
                # increment the count of the ngram
                r.zincrby(SORTED_SET_NAME, 1, ngram)
                # for ZRANGEBYLEX
                r.zadd(LEX_SET_NAME, {ngram: 0})
        print(f"Done processing query \"{query}\"")

def suggest(prefix: str, limit: int = 10, candidate_cap: int = 200) -> list[str]:
        '''
        Suggest a query
        '''
        prefix = prefix.lower()
        if not prefix:
                return []

        # candidates by prefix (alphabetical), capped
        start = f"[{prefix}"
        end = f"[{prefix}\xff"
        candidates = r.zrangebylex(LEX_SET_NAME, start, end, start=0, num=candidate_cap)
        if not candidates:
                return []

        # get scores in one call
        scores = r.zmscore(SORTED_SET_NAME, candidates)

        # sort by score desc, filter Nones to 0, take top K
        paired = [(c, (s or 0.0)) for c, s in zip(candidates, scores)]
        paired.sort(key=lambda x: x[1], reverse=True)
        return [c for c, _ in paired[:limit]]


def main():
        '''
        Main function
        '''
        r.delete(QUEUE_NAME)
        r.delete(SORTED_SET_NAME)
        r.delete(LEX_SET_NAME)
        # --- drain the queue (no blocking) ---
        while True:
                q = dequeue_query()
                if not q:
                        break
                process_dequeued_query(q)



if __name__ == "__main__":
        main()