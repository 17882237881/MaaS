package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	modelpb "maas-platform/shared/proto"
)

// Client wraps the gRPC model client
type Client struct {
	conn   *grpc.ClientConn
	client modelpb.ModelServiceClient
}

// NewClient creates a new gRPC client
func NewClient(address string) (*Client, error) {
	// Set up connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// Connect to server
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	return &Client{
		conn:   conn,
		client: modelpb.NewModelServiceClient(conn),
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CreateModel creates a model via gRPC
func (c *Client) CreateModel(ctx context.Context, req *modelpb.CreateModelRequest) (*modelpb.CreateModelResponse, error) {
	return c.client.CreateModel(ctx, req)
}

// GetModel gets a model via gRPC
func (c *Client) GetModel(ctx context.Context, req *modelpb.GetModelRequest) (*modelpb.GetModelResponse, error) {
	return c.client.GetModel(ctx, req)
}

// ListModels lists models via gRPC
func (c *Client) ListModels(ctx context.Context, req *modelpb.ListModelsRequest) (*modelpb.ListModelsResponse, error) {
	return c.client.ListModels(ctx, req)
}

// UpdateModel updates a model via gRPC
func (c *Client) UpdateModel(ctx context.Context, req *modelpb.UpdateModelRequest) (*modelpb.UpdateModelResponse, error) {
	return c.client.UpdateModel(ctx, req)
}

// DeleteModel deletes a model via gRPC
func (c *Client) DeleteModel(ctx context.Context, req *modelpb.DeleteModelRequest) error {
	_, err := c.client.DeleteModel(ctx, req)
	return err
}

// UpdateModelStatus updates model status via gRPC
func (c *Client) UpdateModelStatus(ctx context.Context, req *modelpb.UpdateModelStatusRequest) (*modelpb.UpdateModelStatusResponse, error) {
	return c.client.UpdateModelStatus(ctx, req)
}

// AddModelTags adds tags to a model via gRPC
func (c *Client) AddModelTags(ctx context.Context, req *modelpb.AddModelTagsRequest) error {
	_, err := c.client.AddModelTags(ctx, req)
	return err
}

// RemoveModelTags removes tags from a model via gRPC
func (c *Client) RemoveModelTags(ctx context.Context, req *modelpb.RemoveModelTagsRequest) error {
	_, err := c.client.RemoveModelTags(ctx, req)
	return err
}

// SetModelMetadata sets model metadata via gRPC
func (c *Client) SetModelMetadata(ctx context.Context, req *modelpb.SetModelMetadataRequest) error {
	_, err := c.client.SetModelMetadata(ctx, req)
	return err
}

// GetModelMetadata gets model metadata via gRPC
func (c *Client) GetModelMetadata(ctx context.Context, req *modelpb.GetModelMetadataRequest) (*modelpb.GetModelMetadataResponse, error) {
	return c.client.GetModelMetadata(ctx, req)
}
