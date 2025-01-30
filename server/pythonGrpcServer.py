import grpc
from concurrent import futures
import torch
from PIL import Image
import io
from transformers import CLIPProcessor, CLIPModel
import embedding_pb2
import embedding_pb2_grpc

class EmbeddingServicer(embedding_pb2_grpc.EmbeddingServiceServicer):
    def __init__(self):
        # CLIPモデルとプロセッサーの初期化
        self.model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")

    def GetEmbedding(self, request, context):
        try:
            # バイトデータからPIL Imageを作成
            image = Image.open(io.BytesIO(request.image_data))
            
            # 画像の前処理とエンベディングの生成
            inputs = self.processor(images=image, return_tensors="pt")
            image_features = self.model.get_image_features(**inputs)
            
            # テンソルをPythonのリストに変換
            embeddings = image_features.detach().numpy().tolist()[0]
            
            return embedding_pb2.EmbeddingResponse(
                success=True,
                embeddings=embeddings
            )
        
        except Exception as e:
            return embedding_pb2.EmbeddingResponse(
                success=False,
                error=str(e)
            )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    embedding_pb2_grpc.add_EmbeddingServiceServicer_to_server(
        EmbeddingServicer(), server
    )
    server.add_insecure_port('[::]:50051')
    print("Starting server on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
