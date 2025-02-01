package util

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "registWebCam/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateClient() (pb.EmbeddingServiceClient, error) {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client := pb.NewEmbeddingServiceClient(conn)

	return client, nil
}

func GetEmbedding(client pb.EmbeddingServiceClient, imageData []byte, filename string) (imgEmbeddng []float32) {

	// コンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// リクエストを作成して送信
	req := &pb.ImageRequest{
		ImageData: imageData,
		Filename:  filename,
	}

	// 開始時間を取得
	startTime := time.Now()
	resp, err := client.GetEmbedding(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get embedding: %v", err)
	}

	endTime := time.Now()
	fmt.Printf("Elapsed time: %v\n", endTime.Sub(startTime))

	if !resp.Success {
		log.Fatalf("Server error: %s", resp.Error)
	}

	fmt.Printf("Successfully got embeddings of length: %d\n", len(resp.Embeddings))
	if len(resp.Embeddings) >= 5 {
		fmt.Println("First 5 embedding values:", resp.Embeddings[:5])
	}

	return resp.Embeddings
}
