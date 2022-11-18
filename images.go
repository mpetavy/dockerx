package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mpetavy/common"
	"github.com/spyzhov/ajson"
)

func listImage(cli *client.Client, image types.ImageSummary, ba []byte) error {
	if *query != "" {
		nodes, err := ajson.JSONPath(ba, *query)
		if err != nil {
			return err
		}

		for _, node := range nodes {
			fmt.Printf("%s\n", TrimApostroph(node.String()))
		}
	} else {
		fmt.Printf("%s\n", string(ba))
	}

	return nil
}

func images(cli *client.Client) error {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if common.Error(err) {
		return err
	}

	for _, image := range images {
		ba, err := json.MarshalIndent(image, "", "    ")
		if common.Error(err) {
			return err
		}

		if *filter != "" {
			nodes, err := ajson.JSONPath(ba, *filter)
			if err != nil {
				return err
			}

			if len(nodes) == 0 {
				continue
			}
		}

		switch {
		case *list:
			err := listImage(cli, image, ba)
			if common.Error(err) {
				return err
			}
		default:
			return &ErrUndefinedAction{}
		}
	}

	return nil
}
