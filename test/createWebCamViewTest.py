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
    # ビューを作成するコマンドを実行
    db.command({
        "create": "webcamViewTest",  # 作成するビュー名
        "viewOn": "webcam",      # 元となるコレクション名
        "pipeline": [
            {
                "$project": {
                    "webcam.webcamid": 1,
                    "webcam.location.latitude": 1,
                    "webcam.location.longitude": 1,
                }
            }
        ]
    })

    print("ビュー 'webcamViewTest' が正常に作成されました。")
    client.close()
except Exception as error:
    print('データベース接続エラー:', error)
    raise

