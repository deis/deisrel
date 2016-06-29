package helm

import (
	"fmt"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
)

// GenParamsComponentAttrs holds the data to change for a single component
type GenParamsComponentAttrs struct {
	Org        string
	Tag        string
	PullPolicy string
}

// GenParamsComponentMap holds the mapping from component to its generate params data
type GenParamsComponentMap map[string]GenParamsComponentAttrs

// CreateGenParamsComponentMap returns an empty GenParamsComponentMap
func CreateGenParamsComponentMap() GenParamsComponentMap {
	return GenParamsComponentMap(map[string]GenParamsComponentAttrs{})
}

func chartPath(chartName string) string {
	return fmt.Sprintf("%s/tpl/generate_params.toml", chartName)
}

func getGenerateParamsTpl(ghClient *github.Client, chartName string) (*template.Template, error) {
	rc, err := ghClient.Repositories.DownloadContents(
		"deis",
		"charts",
		chartPath(chartName),
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		return nil, err
	}
	var paramsData map[string]interface{}
	defer rc.Close()
	md, err := toml.DecodeReader(rc, &paramsData)
	if err != nil {
		return nil, err
	}
	fmt.Println(md.Keys())
	return nil, fmt.Errorf("BOOM")
}

// const (
// 	// TODO: https://github.com/deis/deisrel/issues/11
// 	generateParamsTplStr = `#helm:generate $HELM_GENERATE_DIR/tpl/minio.sh
// #
// # This is the main configuration file for Deis object storage. The values in
// # this file are passed into the appropriate services so that they can configure
// # themselves for persisting data in object storage.
// #
// # In general, all object storage credentials must be able to read and write to
// # the container or bucket they are configured to use.
// #
// # When you change values in this file, make sure to re-run 'helmc generate'
// # on this chart.
//
// # Set the storage backend
// #
// # Valid values are:
// # - filesystem: Store persistent data on ephemeral disk
// # - s3: Store persistent data in AWS S3 (configure in S3 section)
// # - azure: Store persistent data in Azure's object storage
// # - gcs: Store persistent data in Google Cloud Storage
// # - minio: Store persistent data on in-cluster Minio server
// storage = "minio"
//
// [s3]
// accesskey = "YOUR KEY HERE"
// secretkey = "YOUR SECRET HERE"
// # Any S3 region
// region = "us-west-1"
// # Your buckets.
// registry_bucket = "your-registry-bucket-name"
// database_bucket = "your-database-bucket-name"
// builder_bucket = "your-builder-bucket-name"
//
// [azure]
// accountname = "YOUR ACCOUNT NAME"
// accountkey = "YOUR ACCOUNT KEY"
// registry_container = "your-registry-container-name"
// database_container = "your-database-container-name"
// builder_container = "your-builder-container-name"
//
// [gcs]
// # key_json is expanded into a JSON file on the remote server. It must be
// # well-formatted JSON data.
// key_json = '''Paste JSON data here.'''
// registry_bucket = "your-registry-bucket-name"
// database_bucket = "your-database-bucket-name"
// builder_bucket = "your-builder-bucket-name"
//
// [minio]
// org = "{{.Minio.Org}}"
// pullPolicy = "{{.Minio.PullPolicy}}"
// dockerTag = "{{.Minio.Tag}}"
//
// [builder]
// org = "{{.Builder.Org}}"
// pullPolicy = "{{.Builder.PullPolicy}}"
// dockerTag = "{{.Builder.Tag}}"
//
// [slugbuilder]
// org = "{{.SlugBuilder.Org}}"
// pullPolicy = "{{.SlugBuilder.PullPolicy}}"
// dockerTag = "{{.SlugBuilder.Tag}}"
//
// [dockerbuilder]
// org = "{{.DockerBuilder.Org}}"
// pullPolicy = "{{.DockerBuilder.PullPolicy}}"
// dockerTag = "{{.DockerBuilder.Tag}}"
//
// [controller]
// org = "{{.Controller.Org}}"
// pullPolicy = "{{.Controller.PullPolicy}}"
// dockerTag = "{{.Controller.Tag}}"
//
// [slugrunner]
// org = "{{.SlugRunner.Org}}"
// pullPolicy = "{{.SlugRunner.PullPolicy}}"
// dockerTag = "{{.SlugRunner.Tag}}"
//
// [database]
// org = "{{.Database.Org}}"
// pullPolicy = "{{.Database.PullPolicy}}"
// dockerTag = "{{.Database.Tag}}"
//
// [registry]
// org = "{{.Registry.Org}}"
// pullPolicy = "{{.Registry.PullPolicy}}"
// dockerTag = "{{.Registry.Tag}}"
//
// [workflowManager]
// org = "{{.WorkflowManager.Org}}"
// pullPolicy = "{{.WorkflowManager.PullPolicy}}"
// dockerTag = "{{.WorkflowManager.Tag}}"
// versionsApiURL = "https://versions.deis.com"
// doctorApiURL = "https://doctor.deis.com"
//
// [logger]
// org = "{{.Logger.Org}}"
// pullPolicy = "{{.Logger.PullPolicy}}"
// dockerTag = "{{.Logger.Tag}}"
//
// [router]
// org = "{{.Router.Org}}"
// pullPolicy = "{{.Router.PullPolicy}}"
// dockerTag = "{{.Router.Tag}}"
// platformDomain = ""
//
// [fluentd]
// org = "{{.FluentD.Org}}"
// pullPolicy = "{{.FluentD.PullPolicy}}"
// dockerTag = "{{.FluentD.Tag}}"
//
// [grafana]
// org = "{{.Grafana.Org}}"
// pullPolicy = "{{.Grafana.PullPolicy}}"
// dockerTag = "{{.Grafana.Tag}}"
//
// [influxdb]
// org = "{{.InfluxDB.Org}}"
// pullPolicy = "{{.InfluxDB.PullPolicy}}"
// dockerTag = "{{.InfluxDB.Tag}}"
//
// [telegraf]
// org = "{{.Telegraf.Org}}"
// pullPolicy = "{{.Telegraf.PullPolicy}}"
// dockerTag = "{{.Telegraf.Tag}}"
// `
// )
//
// var (
// 	generateParamsTpl = template.Must(template.New("generateParamsTpl").Parse(generateParamsTplStr))
// )
