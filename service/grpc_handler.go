package service

import (
	"context"
	"encoding/json"
	"strings"

	pb "github.com/drone/go-convert/proto"
	"github.com/drone/go-convert/service/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler implements the GoConvertServiceServer interface.
type GRPCHandler struct {
	pb.UnimplementedGoConvertServiceServer
}

func (h *GRPCHandler) ConvertPipeline(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(ctx, entityPipeline, req)
}

func (h *GRPCHandler) ConvertTemplate(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(ctx, entityTemplate, req)
}

func (h *GRPCHandler) ConvertInputSet(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(ctx, entityInputSet, req)
}

func (h *GRPCHandler) ConvertTrigger(ctx context.Context, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	return h.convert(ctx, entityTrigger, req)
}

func (h *GRPCHandler) ConvertExpression(_ context.Context, req *pb.ExpressionConvertRequest) (*pb.ExpressionConvertResponse, error) {
	expr := req.GetExpression()
	exprs := req.GetExpressions()
	remoteFile := req.GetRemoteFile()

	if expr == "" && len(exprs) == 0 && remoteFile == "" {
		return nil, status.Error(codes.InvalidArgument, "one of 'expression', 'expressions', or 'remote_file' field is required")
	}

	// FQN-only context: v1 pipeline YAML plus optional call-site FQN.
	var ctx *converter.ExpressionContext
	if pipelineYAML := strings.TrimSpace(req.GetContextPipelineYaml()); pipelineYAML != "" {
		ctx = &converter.ExpressionContext{
			ContextPipelineYAML: pipelineYAML,
			CurrentFQN:          req.GetCurrentFqn(),
		}
	}

	// Handle remote file
	if remoteFile != "" {
		converted, warnings := converter.ConvertExpressionWithWarnings(remoteFile, ctx)
		return &pb.ExpressionConvertResponse{
			RemoteFile: converted,
			Warnings:   warnings,
			Checksum:   Checksum([]byte(converted)),
		}, nil
	}

	// Handle single expression
	if expr != "" {
		converted, warnings := converter.ConvertExpressionWithWarnings(expr, ctx)
		return &pb.ExpressionConvertResponse{
			Expression: converted,
			Warnings:   warnings,
			Checksum:   Checksum([]byte(converted)),
		}, nil
	}

	// Handle multiple expressions
	converted, warnings := converter.ConvertExpressionsWithWarnings(exprs, ctx)
	return &pb.ExpressionConvertResponse{
		Expressions: converted,
		Warnings:    warnings,
		Checksum:    checksumMap(converted),
	}, nil
}

// checksumMap computes the checksum over the JSON encoding of a map.
func checksumMap(m map[string]string) string {
	b, _ := json.Marshal(m)
	return Checksum(b)
}

func (h *GRPCHandler) GetChecksum(_ context.Context, req *pb.ChecksumRequest) (*pb.ChecksumResponse, error) {
	if strings.TrimSpace(req.GetYaml()) == "" {
		return nil, status.Error(codes.InvalidArgument, "'yaml' field is required and must not be empty")
	}
	return &pb.ChecksumResponse{
		Checksum: Checksum([]byte(req.GetYaml())),
	}, nil
}

func (h *GRPCHandler) convert(ctx context.Context, entityType string, req *pb.ConvertRequest) (*pb.ConvertResponse, error) {
	if strings.TrimSpace(req.GetYaml()) == "" {
		return nil, status.Error(codes.InvalidArgument, "'yaml' field is required and must not be empty")
	}

	// Record entity metadata (account/org/project/id) for the gRPC log line.
	setRequestMetadata(ctx, entityType, req.GetYaml())

	res, err := dispatch(
		entityType,
		req.GetYaml(),
		req.GetTemplateRefMapping(),
		req.GetPipelineRefMapping(),
		req.GetContextPipelineYaml(),
	)
	if err != nil {
		return nil, classifyGRPCError(err)
	}

	return &pb.ConvertResponse{
		Yaml:     string(res.YAML),
		Checksum: Checksum([]byte(req.GetYaml())),
		Report:   toPBReport(buildReport(res.Summary, res.UnknownFields)),
	}, nil
}

// toPBReport converts the public ConversionReport DTO into its proto form.
func toPBReport(r *ConversionReport) *pb.ConversionReport {
	if r == nil {
		return nil
	}
	out := &pb.ConversionReport{
		UnrecognizedFields: r.UnrecognizedFields,
	}
	if len(r.Messages) > 0 {
		out.Messages = make([]*pb.ConverterMessage, 0, len(r.Messages))
		for _, m := range r.Messages {
			out.Messages = append(out.Messages, &pb.ConverterMessage{
				Severity: severityFromString(m.Severity),
				Code:     m.Code,
				Message:  m.Message,
				Context:  m.Context,
			})
		}
	}
	if len(r.Expressions) > 0 {
		out.Expressions = make([]*pb.ExpressionEntry, 0, len(r.Expressions))
		for _, e := range r.Expressions {
			out.Expressions = append(out.Expressions, &pb.ExpressionEntry{
				Original:  e.Original,
				Converted: e.Converted,
				Status:    statusFromString(e.Status),
			})
		}
	}
	return out
}

func severityFromString(s string) pb.Severity {
	switch s {
	case "WARNING":
		return pb.Severity_WARNING
	case "ERROR":
		return pb.Severity_ERROR
	default:
		return pb.Severity_INFO
	}
}

func statusFromString(s string) pb.ConversionStatus {
	if s == "NOT_CONVERTED" {
		return pb.ConversionStatus_NOT_CONVERTED
	}
	return pb.ConversionStatus_SUCCESS
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
