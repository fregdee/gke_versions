package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"context"

	container "cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func main() {
	credentialFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialFilePath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS is empty")
	}
	d, err := os.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatal(err)
	}

	type Credential struct {
		ProjectId string `json:"project_id"`
	}

	var credential Credential

	json.Unmarshal(d, &credential)

	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	parent := fmt.Sprintf("projects/%s/locations/-", credential.ProjectId)

	req := &containerpb.ListClustersRequest{
		Parent: parent,
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cluster Name", "Cluster Version", "NodePool Name", "NodePool Version"})
	table.SetAutoMergeCells(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetRowLine(true)

	var data [][]string
	for _, cluster := range resp.Clusters {
		for _, nodePool := range cluster.NodePools {
			data = append(data, []string{cluster.Name, cluster.CurrentMasterVersion, nodePool.Name, nodePool.Version})
		}
	}

	table.AppendBulk(data)
	table.Render()
}
