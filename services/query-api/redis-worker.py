import redis

r = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)

QUEUE_NAME = "query_queue"

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
                for i in range(len(tokens)):
                        ngram = " ".join(tokens[i:i+n])
                        ngrams.append(ngram)
        return ngrams

def tokenize(query: str) -> list[str]:
    '''
    Tokenize a query
    '''
    return query.lower().split()

def process_query(query: str):
        '''
        Process a query
        '''
        print(f"Processing query: \"{query}\"\n")
        tokenized_query = tokenize(query)
        ngrams = generate_ngrams(tokenized_query)
        print(f"N-grams: {ngrams}\n")
        print(f"Done processing query \"{query}\"")

def main():
        '''
        Main function
        '''
        while True:
                query = dequeue_query()
                if query:
                        process_query(query)

if __name__ == "__main__":
        main()