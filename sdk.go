package sdk

import (
	"context"

	"github.com/go-furnace/proto/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

/*
 *
 * Plugin interface declarations.
 *
 */

// PreCreate is the interface for anything before the build happens. The
// PreCreate plugin has the change to abort the build if returns false.
type PreCreate interface {
	Execute(stackname string) bool
}

// PostCreate interface is the definition of the PostCreate api that can be
// implemented and used via plugins. This interface gives access to the
// stack name.
type PostCreate interface {
	Execute(stackname string)
}

// PreDelete is the interface for anything before the delete happens. The
// PreDelete plugin has the change to abort the build if returns false.
type PreDelete interface {
	Execute(stackname string) bool
}

// PostDelete interface is the definition of the PostDelete api that can be
// implemented and used via plugins. This interface gives access to the
// stack name.
type PostDelete interface {
	Execute(stackname string)
}

/*
 *
 * PRECREATE Plugin structs and functions.
 *
 */

// PreCreateGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type PreCreateGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PreCreate
}

// GRPCServer is the grpc server implementation which calls the
// protoc generated code to register it.
func (p *PreCreateGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPreCreateServer(s, &GRPCPreCreateServer{Impl: p.Impl})
	return nil
}

// GRPCClient is the grpc client that will talk to the GRPC Server
// and calls into the generated protoc code.
func (p *PreCreateGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCPreCreateClient{client: proto.NewPreCreateClient(c)}, nil
}

// GRPCPreCreateClient is an implementation of PreCreate that talks over RPC.
type GRPCPreCreateClient struct{ client proto.PreCreateClient }

// Execute is the GRPC implementation of the Execute function for the
// PreCreate plugin definition. This will talk over GRPC.
func (m *GRPCPreCreateClient) Execute(key string) bool {
	p, err := m.client.Execute(context.Background(), &proto.Stack{
		Name: key,
	})
	if err != nil {
		return false
	}
	return p.Failed
}

// GRPCPreCreateServer is the gRPC server that GRPCPreCreateClient talks to.
type GRPCPreCreateServer struct {
	// This is the real implementation
	Impl PreCreate
}

// Execute is the execute function of the GRPCServer which will rely the information to the
// underlying implementation of this interface.
func (m *GRPCPreCreateServer) Execute(ctx context.Context, req *proto.Stack) (*proto.Proceed, error) {
	res := m.Impl.Execute(req.Name)
	return &proto.Proceed{Failed: res}, nil
}

/*
 *
 * POSTCREATE Plugin structs and functions.
 *
 */

// PostCreateGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type PostCreateGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PostCreate
}

// GRPCPostCreateClient is an implementation of PreCreate that talks over RPC.
type GRPCPostCreateClient struct{ client proto.PostCreateClient }

// Execute is the GRPC implementation of the Execute function for the
// PostCreate plugin definition. This will talk over GRPC.
func (m *GRPCPostCreateClient) Execute(stackname string) {
	m.client.Execute(context.Background(), &proto.Stack{
		Name: stackname,
	})
}

// GRPCPostCreateServer is the gRPC server that GRPCPostCreateClient talks to.
type GRPCPostCreateServer struct {
	// This is the real implementation
	Impl PostCreate
}

// GRPCServer is the grpc server implementation which calls the
// protoc generated code to register it.
func (p *PostCreateGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPostCreateServer(s, &GRPCPostCreateServer{Impl: p.Impl})
	return nil
}

// GRPCClient is the grpc client that will talk to the GRPC Server
// and calls into the generated protoc code.
func (p *PostCreateGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCPostCreateClient{client: proto.NewPostCreateClient(c)}, nil
}

// Execute is the execute functin of the GRPCServer which will rely the information to the
// underlying implementation of this interface.
func (m *GRPCPostCreateServer) Execute(ctx context.Context, req *proto.Stack) (*proto.Empty, error) {
	m.Impl.Execute(req.Name)
	return &proto.Empty{}, nil
}

/*
 *
 * PREDELETE Plugin structs and functions.
 *
 */

// PreDeleteGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type PreDeleteGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PreDelete
}

// GRPCServer is the grpc server implementation which calls the
// protoc generated code to register it.
func (p *PreDeleteGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPreDeleteServer(s, &GRPCPreDeleteServer{Impl: p.Impl})
	return nil
}

// GRPCClient is the grpc client that will talk to the GRPC Server
// and calls into the generated protoc code.
func (p *PreDeleteGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCPreDeleteClient{client: proto.NewPreDeleteClient(c)}, nil
}

// GRPCPreDeleteClient is an implementation of PreDelete that talks over RPC.
type GRPCPreDeleteClient struct{ client proto.PreDeleteClient }

// Execute is the GRPC implementation of the Execute function for the
// PreDelete plugin definition. This will talk over GRPC.
func (m *GRPCPreDeleteClient) Execute(key string) bool {
	p, err := m.client.Execute(context.Background(), &proto.Stack{
		Name: key,
	})
	if err != nil {
		return false
	}
	return p.Failed
}

// GRPCPreDeleteServer is the gRPC server that GRPCPreCreateClient talks to.
type GRPCPreDeleteServer struct {
	// This is the real implementation
	Impl PreDelete
}

// Execute is the execute function of the GRPCServer which will rely the information to the
// underlying implementation of this interface.
func (m *GRPCPreDeleteServer) Execute(ctx context.Context, req *proto.Stack) (*proto.Proceed, error) {
	res := m.Impl.Execute(req.Name)
	return &proto.Proceed{Failed: res}, nil
}

/*
 *
 * POSTCREATE Plugin structs and functions.
 *
 */

// PostDeleteGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type PostDeleteGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl PostDelete
}

// GRPCPostDeleteClient is an implementation of PreCreate that talks over RPC.
type GRPCPostDeleteClient struct{ client proto.PostDeleteClient }

// Execute is the GRPC implementation of the Execute function for the
// PostDelete plugin definition. This will talk over GRPC.
func (m *GRPCPostDeleteClient) Execute(stackname string) {
	m.client.Execute(context.Background(), &proto.Stack{
		Name: stackname,
	})
}

// GRPCPostDeleteServer is the gRPC server that GRPCPostDeleteClient talks to.
type GRPCPostDeleteServer struct {
	// This is the real implementation
	Impl PostDelete
}

// GRPCServer is the grpc server implementation which calls the
// protoc generated code to register it.
func (p *PostDeleteGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPostDeleteServer(s, &GRPCPostDeleteServer{Impl: p.Impl})
	return nil
}

// GRPCClient is the grpc client that will talk to the GRPC Server
// and calls into the generated protoc code.
func (p *PostDeleteGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCPostDeleteClient{client: proto.NewPostDeleteClient(c)}, nil
}

// Execute is the execute functin of the GRPCServer which will rely the information to the
// underlying implementation of this interface.
func (m *GRPCPostDeleteServer) Execute(ctx context.Context, req *proto.Stack) (*proto.Empty, error) {
	m.Impl.Execute(req.Name)
	return &proto.Empty{}, nil
}
