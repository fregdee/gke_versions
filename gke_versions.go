package gke_versions

import (
	"context"
	"encoding/json"
	//"errors"
	"github.com/pkg/errors"
	"fmt"
	"os"
	"runtime"

	"github.com/olekukonko/tablewriter"
	"github.com/jessevdk/go-flags"

	container "cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

const cmdName = "gke_versions"

// version by Makefile
var version string

// Required permissions
// - container.clusters.list
const googleApplicationCredentialsPath = "GOOGLE_APPLICATION_CREDENTIALS"

type cmdOpts struct {
	Version bool `short:"v" long:"version" description:"show version"`
}

func printVersion() {
	fmt.Printf(`%s %s
Compiler: %s %s
`,
		cmdName,
		version,
		runtime.Compiler,
		runtime.Version())
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

func Run() error {
	opts := cmdOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	psr.Name = cmdName
	_, err := psr.Parse()
	if err != nil {
		// TODO: handle err
		return nil
	}

	if opts.Version {
		printVersion()
		return nil
	}

	projectName, err := getProjectName()
	if err != nil {
		return err
	}

	req := &containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/-", projectName),
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
		fmt.Println("0 clusters")
		return nil
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
