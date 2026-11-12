package aikit

import (
	"context"
	"errors"
	"os"

	"github.com/firebase/genkit/go/genkit"
	genkitmcp "github.com/firebase/genkit/go/plugins/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

const mcpServerName string = "Stan"

type MCPServer struct {
	*genkitmcp.GenkitMCPServer
}

func NewMCPServer(ctx context.Context, tools []Tool) *MCPServer {
	g := genkit.Init(ctx)

	for _, tool := range tools {
		genkit.RegisterAction(g, tool.Ref)
	}

	return &MCPServer{
		GenkitMCPServer: genkitmcp.NewMCPServer(
			g,
			genkitmcp.MCPServerOptions{
				Name: mcpServerName,
			}),
	}
}

func (s *MCPServer) Start(ctx context.Context) error {
	if err := mcpserver.NewStdioServer(s.GetServer()).Listen(ctx, os.Stdin, os.Stdout); err != nil && !errors.Is(
		err,
		context.Canceled,
	) {
		return err
	}

	return nil
}
