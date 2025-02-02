from pymongo import MongoClient

uri = 'mongodb://localhost:27017'
client = MongoClient(uri)

def connect_to_database():
    try:
        # データベース接続を確認
        client.admin.command('ismaster')
        db = client['your_database_name']
        return db
    except Exception as error:
        print('データベース接続エラー:', error)
        raise

try:
    db = connect_to_database()
    db.webcam.aggregate([
    {
        "$vectorSearch": {
            "index": "imgembindex",
            "path": "embedding",
            "queryVector": [<array-of-numbers>],
            "numCandidates": 100,
            "limit": 10
        }
    }
    ])
except Exception as error:
    print('データベース接続エラー:', error)
    raise

