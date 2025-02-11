import requests
import dotenv
import os

# GraphQLエンドポイントURL
url = "https://webcam.hasura.app/v1/graphql"

# GraphQLクエリ
query = """
query {
  webcamView(limit: 200, order_by: {webcam: {webcamid: asc}}) {
    webcam {
      location {
        city
        country
        latitude
        longitude
        region
      }
      title
      webcamid
    }
  }
}
"""

dotenv.load_dotenv()

secret = os.environ.get("HASURA_ADMINE_SECRET")
# ヘッダー (必要に応じて認証トークンを追加)
headers = {
    "Content-Type": "application/json",
    "x-hasura-admin-secret": secret ,   # 認証が必要な場合に追加
}

# リクエストの送信
response = requests.post(url, json={"query": query}, headers=headers)

# レスポンスの確認と出力
if response.status_code == 200:
    print("Query successful!")
    print(response.json())  # JSONデータを表示
else:
    print(f"Query failed with status code {response.status_code}")
    print(response.text)  # エラーメッセージを表示
