package rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/sagacious-labs/k8trics/pkg/protos/v1alpha1/api"
	"github.com/sagacious-labs/k8trics/pkg/store"
	"github.com/sagacious-labs/k8trics/pkg/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// WatchDataResponse represents the response of the watch RPCs
type WatchDataResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
}

// HyperionApply is a wrapper around hyperion's `Apply` RPC
func HyperionApply(ctx context.Context, req *api.ApplyRequest, host string) (*api.ApplyResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
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
func HyperionDelete(ctx context.Context, req *api.DeleteRequest, host string) (*api.DeleteResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
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
func HyperionGet(ctx context.Context, req *api.GetRequest, host string) (*api.GetResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
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
func HyperionList(ctx context.Context, req *api.ListRequest, host string) (chan *api.GetResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
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
				return
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
func HyperionWatchData(ctx context.Context, req *api.WatchDataRequest, host string) (chan *WatchDataResponse, error) {
	store, ok := ctx.Value("pod_store").(*store.PodStore)
	if !ok {
		return nil, errors.New("pod store not found")
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.WatchData(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *WatchDataResponse, 8)

	go func() {
		for {
			item, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					close(ch)
				}

				return
			}

			data := parseWatchDataJSON(item.Data)
			cid, ok := data["container_id"].(string)
			if !ok {
				logrus.Warn("container_id not found in the retrieved data")
				continue
			}

			pod, ok := store.GetByContainerID(cid)
			if !ok {
				logrus.Warn("no pod found for container id: ", cid)
				continue
			}
			data["name"] = utils.TrimPodTemplateHash(&pod.Pod)

			ch <- &WatchDataResponse{
				Data: data,
			}
		}
	}()

	return ch, nil
}

// HyperionWatchLog is a wrapper around hyperion's `WatchLog` RPC
func HyperionWatchLog(ctx context.Context, req *api.WatchLogRequest, host string) (chan string, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewHyperionAPIServiceClient(conn)
	res, err := client.WatchLog(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan string, 8)

	go func() {
		for {
			item, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					close(ch)
				}

				return
			}

			ch <- parseWatchLog(item.Data)
		}
	}()

	return ch, nil
}

// parseWatchDataJSON takes in a slice of byte and converts it into
// map[string]interface{}
//
// If the method fails then it returns an empty map
func parseWatchDataJSON(byt []byte) (mp map[string]interface{}) {
	_ = json.Unmarshal(byt, &mp)
	return
}

func parseWatchLog(byt []byte) string {
	res := []byte{}

	_, err := base64.StdEncoding.Decode(res, byt)
	if err != nil {
		logrus.Warn("failed to decode the log message")
	}

	return string(res)
}
