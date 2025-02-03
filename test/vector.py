from pymongo import MongoClient
from transformers import AutoProcessor,CLIPProcessor, CLIPModel, CLIPTokenizer
from dotenv import load_dotenv
from os import getenv

load_dotenv()

uri = getenv("MONGODB_URI") 
client = MongoClient(uri)

def connect_to_database():
    try:
        # データベース接続を確認
        client.admin.command('ismaster')
        db = client['webcamNew']
        return db
    except Exception as error:
        print('データベース接続エラー:', error)
        raise


def get_model_info(model_ID, device):
	model = CLIPModel.from_pretrained(model_ID).to(device)
	processor = AutoProcessor.from_pretrained(model_ID)
	tokenizer = CLIPTokenizer.from_pretrained(model_ID)
    # Return model, processor & tokenizer
	return model, processor, tokenizer

# Set the device
#device = "cuda" if torch.cuda.is_available() else "cpu"
device = "cpu"
model_ID = "openai/clip-vit-base-patch32"

model, processor, tokenizer = get_model_info(model_ID, device)

def get_single_text_embedding(text): 
    inputs = tokenizer(text, return_tensors = "pt")
    # normalize input embeddings
    text_embeddings = model.get_text_features(**inputs)
    text_embeddings /= text_embeddings.norm(dim=-1, keepdim=True)     
    # convert the embeddings to numpy array
    return text_embeddings.cpu().detach().numpy()


def get_single_image_embedding(my_image):
    image = processor(images=my_image , return_tensors="pt")
    embedding = model.get_image_features(**image).float()
    # convert the embeddings to numpy array
    return embedding.cpu().detach().numpy()

query ="big city"

vector = get_single_text_embedding(query)
list = vector.tolist()

try:
    db = connect_to_database()
    results = db.webcam.aggregate([
        {
            "$vectorSearch": {
                "index": "imgembindex",
                "path": "embedding",
                "queryVector": list[0],
                "numCandidates": 100,
                "limit": 10
            }
        },
        {
            "$project": {
                "score": { "$meta": "vectorSearchScore" }, 
                "webcam.webcamid": 1
            }
        }
    ])
    fileserver = getenv("FILE_SERVER")
    for result in results:
        webcamid = result.get('webcam').get('webcamid')
        downloadUrl = fileserver + str(webcamid) +".jpg"
        print(downloadUrl)
except Exception as error:
    print('データベース接続エラー:', error)
    raise

