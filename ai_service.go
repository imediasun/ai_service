package main

import (
    "context"
    "log"
    "net"
    "math/rand"
    "google.golang.org/grpc"
    pb "generated/ai_service"
)

// AIService implements the gRPC service for adjusting team strengths
type AIService struct {
    pb.UnimplementedAIServiceServer
}

// AdjustTeamStrength adjusts the strength of teams based on their performance over time
func (s *AIService) AdjustTeamStrength(ctx context.Context, req *pb.TeamStrengthRequest) (*pb.TeamStrengthResponse, error) {
    response := &pb.TeamStrengthResponse{}

    for _, team := range req.Teams {
        // Increase the strength of weaker teams over time with some random factor
        if team.Strength < 3 {
            team.Strength += float32(rand.Float32() * 0.5) // Example of adjusting strength
        }
        response.Teams = append(response.Teams, team)
    }

    return response, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterAIServiceServer(s, &AIService{})
    log.Printf("gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
