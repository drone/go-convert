package service

import (
	"context"
	"strings"

	pb "github.com/drone/go-convert/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler implements the GoConvertServiceServer interface.
type GRPCHandler struct {
	pb.UnimplementedGoConvertServiceServer
}

func (h *GRPCHandler) ConvertPipeline(_ context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(entityPipeline, req)
}

func (h *GRPCHandler) ConvertTemplate(_ context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(entityTemplate, req)
}

func (h *GRPCHandler) ConvertInputSet(_ context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(entityInputSet, req)
}

func (h *GRPCHandler) GetChecksum(_ context.Context, req *pb.ChecksumRequest) (*pb.ChecksumResponse, error) {
	if strings.TrimSpace(req.GetYaml()) == "" {
		return nil, status.Error(codes.InvalidArgument, "'yaml' field is required and must not be empty")
	}
	return &pb.ChecksumResponse{
		Checksum: Checksum([]byte(req.GetYaml())),
	}, nil
}

func (h *GRPCHandler) convert(entityType string, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	if strings.TrimSpace(req.GetYaml()) == "" {
		return nil, status.Error(codes.InvalidArgument, "'yaml' field is required and must not be empty")
	}

	refMapping := mergeRefMappings(req.GetTemplateRefMapping(), req.GetPipelineRefMapping())

	outBytes, err := dispatch(entityType, req.GetYaml(), refMapping)
	if err != nil {
		return nil, classifyGRPCError(err)
	}

	return &pb.ConvertResponse{
		Yaml:     string(outBytes),
		Checksum: Checksum([]byte(req.GetYaml())),
	}, nil
}

func mergeRefMappings(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func classifyGRPCError(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "expected top-level"):
		return status.Error(codes.InvalidArgument, msg)
	case strings.Contains(msg, "failed to parse"):
		return status.Error(codes.InvalidArgument, msg)
	case strings.Contains(msg, "unsupported template type"),
		strings.Contains(msg, "conversion returned nil"),
		strings.Contains(msg, "produced no stages"),
		strings.Contains(msg, "produced no steps"):
		return status.Error(codes.InvalidArgument, msg)
	case strings.Contains(msg, "unknown entity_type"):
		return status.Error(codes.InvalidArgument, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}
