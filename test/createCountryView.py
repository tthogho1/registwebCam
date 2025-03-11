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
        db = client['webcam']        
        return db
    except Exception as error:
        print('データベース接続エラー:', error)
        raise

try:
    db = connect_to_database()
    # ビューを作成するコマンドを実行
    db.command({
        "create": "Country",
        "viewOn": "webcam",
        "pipeline": [
            {
                "$group": {
                    "_id": "$location.country",
                    "code": { "$first": "$location.countrycode" },
                    "country": { "$first": "$location.country" }
                }
            },
            {
                "$project": {
                    "_id": 0,
                    "code": 1,
                    "country": 1
                }
            }
        ]
    })

    print("ビュー 'CountryView' が正常に作成されました。")
    
except Exception as error:
    print('データベース接続エラー:', error)
    raise

