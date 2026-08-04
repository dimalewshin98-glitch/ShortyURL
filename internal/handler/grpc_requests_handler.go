package handler

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/dimalewshin98-glitch/ShortyURL/api"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/repository"
	"github.com/dimalewshin98-glitch/ShortyURL/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCRequestsHandler struct {
	pb.UnimplementedShortenerServiceServer
	shortenerService service.ServiceInterface
}

func NewGRPCRequestsHandler(service service.ServiceInterface) *GRPCRequestsHandler {
	return &GRPCRequestsHandler{
		shortenerService: service,
	}
}

func (g *GRPCRequestsHandler) ShortenURL(ctx context.Context, in *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	rawUserID := ctx.Value("userID")
	if rawUserID == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}
	userID, ok := rawUserID.(int)
	if !ok {
		return nil, fmt.Errorf("invalid user ID type in context")
	}
	url := in.GetUrl()
	if url == "" {
		return nil, fmt.Errorf("URL is empty")
	}
	shortURL, err := g.shortenerService.Shorten(ctx, userID, url)
	if err != nil {
		return nil, err
	}
	response := &pb.URLShortenResponse_builder{
		Result: &shortURL,
	}
	return response.Build(), nil
}

func (g *GRPCRequestsHandler) ExpandURL(ctx context.Context, in *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	rawUserID := ctx.Value("userID")
	if rawUserID == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}
	userID, ok := rawUserID.(int)
	if !ok {
		return nil, fmt.Errorf("invalid user ID type in context")
	}
	urlID := in.GetId()
	if urlID == "" {
		return nil, fmt.Errorf("URL ID is empty")
	}
	originalURL, err := g.shortenerService.GetURL(ctx, userID, urlID)
	if err != nil {
		if errors.Is(err, repository.ErrShortURLDeleted) {
			return nil, fmt.Errorf("URL has been deleted: %w", err)
		}
		return nil, err
	}
	response := &pb.URLExpandResponse_builder{
		Result: &originalURL,
	}
	return response.Build(), nil
}

func (g *GRPCRequestsHandler) ListUserURLs(ctx context.Context, in *emptypb.Empty) (*pb.UserURLsResponse, error) {
	rawUserID := ctx.Value("userID")
	if rawUserID == nil {
		return nil, fmt.Errorf("user ID not found in context")
	}
	userID, ok := rawUserID.(int)
	if !ok {
		return nil, fmt.Errorf("invalid user ID type in context")
	}
	userURLs, err := g.shortenerService.UserUrls(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userURLs == nil {
		result := &pb.UserURLsResponse_builder{
			Url: []*pb.URLData{},
		}
		return result.Build(), nil
	}
	var pbURLs []*pb.URLData
	for _, url := range userURLs {
		pbURL := &pb.URLData_builder{
			ShortUrl:    &url.ShortURL,
			OriginalUrl: &url.OriginalURL,
		}
		pbURLs = append(pbURLs, pbURL.Build())
	}
	response := &pb.UserURLsResponse_builder{
		Url: pbURLs,
	}
	return response.Build(), nil
}
