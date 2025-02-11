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
        "create": "webcamView",  # 作成するビュー名
        "viewOn": "webcam",      # 元となるコレクション名
        "pipeline": [
            {
                "$project": {
                    "webcam.title": 1,
                    "webcam.webcamid": {
                        "$convert": {
                            "input": "$webcam.webcamid",  # 元の値
                            "to": "string",               # 変換先データ型
                        }
                    },
                    "webcam.location.city": 1,
                    "webcam.location.region": 1,
                    "webcam.location.country": 1,
                    "webcam.location.latitude": {
                        "$convert": {
                            "input": "$webcam.location.latitude",  # 元の値
                            "to": "string",               # 変換先データ型
                        }
                    },
                    "webcam.location.longitude": {
                        "$convert": {
                            "input": "$webcam.location.longitude",  # 元の値
                            "to": "string",               # 変換先データ型
                        }
                    },
                    "webcam.player.day": 1
                }
            }
        ]
    })

    print("ビュー 'webcamView' が正常に作成されました。")
    
except Exception as error:
    print('データベース接続エラー:', error)
    raise

