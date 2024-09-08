package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "path/to/your/proto/generated/ai_service"
)

type predictionService struct {
    pb.UnimplementedPredictionServiceServer
}

// Функция для получения предсказаний
func (s *predictionService) GetChampionshipPredictions(ctx context.Context, req *pb.TeamsRequest) (*pb.PredictionsResponse, error) {
    var totalPoints int32 = 0

    // Суммируем очки всех команд
    for _, team := range req.Teams {
        totalPoints += team.Points
    }

    predictions := make([]*pb.Prediction, len(req.Teams))

    // Если у всех команд 0 очков, возвращаем предсказания по 0%
    if totalPoints == 0 {
        for i, team := range req.Teams {
            predictions[i] = &pb.Prediction{
                Team:       team.Name,
                Prediction: "0%",
            }
        }
    } else {
        // Вычисляем предсказания на основе очков каждой команды
        for i, team := range req.Teams {
            percentage := float64(team.Points) / float64(totalPoints) * 100
            predictions[i] = &pb.Prediction{
                Team:       team.Name,
                Prediction: fmt.Sprintf("%.2f%%", percentage),
            }
        }
    }

    return &pb.PredictionsResponse{
        Predictions: predictions,
    }, nil
}

func main() {
    lis, err := net.Listen("tcp", "0.0.0.0:50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterPredictionServiceServer(grpcServer, &server{})

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}