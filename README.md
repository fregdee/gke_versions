# gke_versions

## Description
gke_versions command prints versions of all GKE Clusters and NodePool

### Required GCP IAM permission
`container.clusters.list`
See also: https://cloud.google.com/iam/docs/permissions-reference#container.clusters.list

### Authorization
Use [Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)

## Installation
```shell
go get github.com/fregdee/gke_versoins
```

## Usage
```shell
$ gke_versions
+-----------------------------------------+-----------------+---------------------+------------------+
|              CLUSTER NAME               | CLUSTER VERSION |    NODEPOOL NAME    | NODEPOOL VERSION |
+-----------------------------------------+-----------------+---------------------+------------------+
| cluster1                                | 1.17.15-gke.800 | default-pool        | 1.17.13-gke.2001 |
+-----------------------------------------+                 +---------------------+------------------+
| cluster2                                |                 | default-pool        | 1.17.15-gke.800  |
+-----------------------------------------+                 +---------------------+------------------+
| cluster3                                |                 | pool-20201005       | 1.17.12-gke.2502 |
+                                         +                 +---------------------+------------------+
|                                         |                 | monitoring-20201005 | 1.17.12-gke.2502 |
+-----------------------------------------+                 +---------------------+------------------+
| cluster4                                |                 | default-pool        | 1.17.9-gke.6300  |
+-----------------------------------------+-----------------+---------------------+------------------+
```
