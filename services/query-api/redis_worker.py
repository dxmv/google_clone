import redis
import time
import os

# Get Redis connection details from environment variables
redis_host = os.getenv('REDIS_HOST', 'localhost')
redis_port = int(os.getenv('REDIS_PORT', 6379))

r = redis.Redis(host=redis_host, port=redis_port, db=0, decode_responses=True)

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
        Generate n-grams from a query for n=1,2,3
        '''
        ngrams = []
        for n in [1,2,3]:
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
        print(f"-----------------\nProcessing query: \"{query}\"")
        tokenized_query = tokenize(query)
        ngrams = generate_ngrams(tokenized_query)
        print(f"Done generating n-grams for query: \"{query}\"")
        # fill redis sorted set with ngrams
        for ngram in ngrams:
                # increment the count of the ngram
                r.zincrby(SORTED_SET_NAME, 1, ngram)
                # for ZRANGEBYLEX
                r.zadd(LEX_SET_NAME, {ngram: 0})
        print(f"Done processing query \"{query}\"\n-----------------")

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
        # seed the queue with some queries
        enqueue_query("mathematics")
        enqueue_query("math")
        enqueue_query("maths")
        enqueue_query("mathematics")
        enqueue_query("investing")
        enqueue_query("invest")
        enqueue_query("investing in stocks")
        enqueue_query("investing in stocks and bonds")
        enqueue_query("investing in stocks and bonds and etfs")
        enqueue_query("investing in stocks and bonds and etfs and mutual funds")
        enqueue_query("invest")
        # --- drain the queue (no blocking) ---
        while True:
                q = dequeue_query()
                if not q:
                        time.sleep(1)
                        continue
                process_dequeued_query(q)



if __name__ == "__main__":
        main()