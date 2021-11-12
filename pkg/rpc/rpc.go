package rpc

import (
	"context"
	"fmt"
	"io"

	"github.com/sagacious-labs/k8trics/pkg/protos/v1alpha1/api"
	"github.com/sagacious-labs/k8trics/pkg/utils"
	"google.golang.org/grpc"
)

func host() string {
	return fmt.Sprintf("%s:%s", utils.GetEnv("HYPERION_HOST", "0.0.0.0"), utils.GetEnv("HYPERION_PORT", "2310"))
}

// HyperionApply is a wrapper around hyperion's `Apply` RPC
func HyperionApply(ctx context.Context, req *api.ApplyRequest) (*api.ApplyResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.Apply(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HyperionDelete is a wrapper around hyperion's `Delete` RPC
func HyperionDelete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HyperionGet is a wrapper around hyperion's `Get` RPC
func HyperionGet(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HyperionList is a wrapper around hyperion's `List` RPC
func HyperionList(ctx context.Context, req *api.ListRequest) (chan *api.GetResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.List(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *api.GetResponse, 8)

	go func() {
		for {
			item, err := res.Recv()
			if err == io.EOF {
				close(ch)
			}
			if err != nil {
				fmt.Println("[List Error]: ", err)
				continue
			}

			ch <- item
		}
	}()

	return ch, nil
}

// HyperionWatchData is a wrapper around hyperion's `WatchData` RPC
func HyperionWatchData(ctx context.Context, req *api.WatchDataRequest) (chan *api.WatchDataResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.WatchData(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *api.WatchDataResponse, 8)

	go func() {
		for {
			item, err := res.Recv()
			if err == io.EOF {
				close(ch)
			}
			if err != nil {
				continue
			}

			ch <- item
		}
	}()

	return ch, nil
}

// HyperionWatchLog is a wrapper around hyperion's `WatchLog` RPC
func HyperionWatchLog(ctx context.Context, req *api.WatchLogRequest) (chan *api.WatchLogResponse, error) {
	conn, err := grpc.Dial(host(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.WatchLog(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *api.WatchLogResponse, 8)

	go func() {
		for {
			item, err := res.Recv()
			if err == io.EOF {
				close(ch)
			}
			if err != nil {
				continue
			}

			ch <- item
		}
	}()

	return ch, nil
}