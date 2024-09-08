package main

import (
    "context"
    "fmt"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/imediasun/ai_service/generated"
)

type predictionService struct {
    pb.UnimplementedPredictionServiceServer
}

// Function for getting predictions
func (s *predictionService) GetChampionshipPredictions(ctx context.Context, req *pb.TeamsRequest) (*pb.PredictionsResponse, error) {
    var totalPoints int32 = 0

    // We sum up the points of all teams
    for _, team := range req.Teams {
        totalPoints += team.Points
    }

    predictions := make([]*pb.Prediction, len(req.Teams))

    // If all teams have 0 points, we return predictions at 0%
    if totalPoints == 0 {
        for i, team := range req.Teams {
            predictions[i] = &pb.Prediction{
                Team:       team.Name,
                Prediction: "0%",
            }
        }
    } else {
        // We calculate predictions based on the points of each team
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
    pb.RegisterPredictionServiceServer(grpcServer, &predictionService{})

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}