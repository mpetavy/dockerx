package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mpetavy/common"
	"github.com/spyzhov/ajson"
	"time"
)

func listContainer(cli *client.Client, container types.Container, ba []byte) error {
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

func killContainer(cli *client.Client, container types.Container) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Printf("%s\n", container.ID)

	err := cli.ContainerKill(ctx, container.ID, "")
	if common.Error(err) {
		return err
	}

	if err == nil && ctx.Err() != nil {
		err = ctx.Err()
	}

	return nil
}

func containers(cli *client.Client) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if common.Error(err) {
		return err
	}

	for _, container := range containers {
		ba, err := json.MarshalIndent(container, "", "    ")
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
			err := listContainer(cli, container, ba)
			if common.Error(err) {
				return err
			}
		case *kill:
			err := killContainer(cli, container)
			if common.Error(err) {
				return err
			}
		default:
			return &ErrUndefinedAction{}
		}
	}

	return nil
}
