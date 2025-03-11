from pymongo.mongo_client import MongoClient
from pymongo.operations import SearchIndexModel
import time
from dotenv import load_dotenv
import os

load_dotenv()
# Connect to your Atlas deployment
uri = os.getenv('MONGODB_URI')
client = MongoClient(uri)

# Access your database and collection
database = client["webcamNew"]
collection = database["webcam"]

# Create your index model, then create the search index
search_index_model = SearchIndexModel(
  definition={
    "fields": [
      {
        "type": "vector",
        "path": "embedding",
        "numDimensions": 512,
        "similarity": "cosine",
        "quantization": "scalar"
      }
    ]
  },
  name="imgembindex",
  type="vectorSearch",
)

result = collection.create_search_index(model=search_index_model)
print("New search index named " + result + " is building.")

# Wait for initial sync to complete
print("Polling to check if the index is ready. This may take up to a minute.")
predicate=None
if predicate is None:
  predicate = lambda index: index.get("queryable") is True

while True:
  indices = list(collection.list_search_indexes(result))
  if len(indices) and predicate(indices[0]):
    break
  time.sleep(5)
print(result + " is ready for querying.")

client.close()
