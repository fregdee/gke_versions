package gke_versions

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"

	container "cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

const cmdName = "gke_versions"

// Required permissions
// - container.clusters.list
const googleApplicationCredentialsPath = "GOOGLE_APPLICATION_CREDENTIALS"

func Run(outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	flag.Usage = func() {
		fmt.Fprintf(outStream,"%s prints versions of all GKE Clusters and NodePool\n\n")
		fmt.Fprintf(outStream,"Usage:\n\t %s\n", cmdName)
		flag.PrintDefaults()
	}

	projectName, err := getProjectName()
	if err != nil {
		return err
	}

	parent := fmt.Sprintf("projects/%s/locations/-", projectName)

	req := &containerpb.ListClustersRequest{
		Parent: parent,
	}

	ctx := context.Background()
	c, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	resp, err := c.ListClusters(ctx, req)
	if err != nil {
		return err
	}

	clusters := resp.Clusters
	if len(clusters) == 0 {
		_, err := fmt.Fprint(outStream, "0 clusters\n")
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cluster Name", "Cluster Version", "NodePool Name", "NodePool Version"})
	table.SetAutoMergeCells(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetRowLine(true)

	var data [][]string
	for _, cluster := range clusters {
		for _, nodePool := range cluster.NodePools {
			data = append(data, []string{cluster.Name, cluster.CurrentMasterVersion, nodePool.Name, nodePool.Version})
		}
	}

	table.AppendBulk(data)
	table.Render()

	return nil
}

func getProjectName() (string, error) {
	credentialFilePath := os.Getenv(googleApplicationCredentialsPath)
	if credentialFilePath == "" {
		return "", errors.New("GOOGLE_APPLICATION_CREDENTIALS is empty.\nSee also https://cloud.google.com/docs/authentication/getting-started#setting_the_environment_variable")
	}
	d, err := os.ReadFile(credentialFilePath)
	if err != nil {
		return "", err
	}

	type Credential struct {
		ProjectId string `json:"project_id"`
	}

	var credential Credential

	err = json.Unmarshal(d, &credential)
	if err != nil {
		return "", err
	}

	return credential.ProjectId, nil
}
