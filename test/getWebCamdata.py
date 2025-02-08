from pymongo import MongoClient
from dotenv import load_dotenv
from os import getenv
import json

load_dotenv()

uri = getenv("MONGODB_URI") 
client = MongoClient(uri)

def connect_to_database():
    try:
        # データベース接続を確認
        #client.admin.command('ismaster')
        db = client['webcamNew']        
        return db
    except Exception as error:
        print('データベース接続エラー:', error)
        raise

try:
    db = connect_to_database()
    #resutlt = db.webcam.find_one()
    cursor = db.webcam.find()
    for doc in cursor:
        print(json.dumps(doc, indent=4, default=str))  # JSON形式で整形して
    
except Exception as error:
    print('データベース接続エラー:', error)
    raise

