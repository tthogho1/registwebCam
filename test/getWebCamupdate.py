from pymongo import MongoClient
from dotenv import load_dotenv
from os import getenv

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
    # updateManyクエリ
    result = db.webcam.update_many(
        {"$numberInt": {"$type": "string"}},  # 対象フィールドが文字列の場合
        [
            {"$set": {"numberInt": {"$toInt": "$numberInt"}}}  # double型に変換
        ]
    )

    print(result)
    
except Exception as error:
    print('データベース接続エラー:', error)
    raise

