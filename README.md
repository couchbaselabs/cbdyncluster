This is a CLI tool for interfacing with a running cbdynclusterd.

Download the latest release for your platform from https://github.com/couchbaselabs/cbdyncluster/releases/latest

Example Usage (as of May 23, 2022):

Init config

```bash
> cbdyncluster init --auth jake.rawsthorne@couchbase.com --server 172.23.104.149:19923
> cat $HOME/.cbdyncluster.toml
auth = "jake.rawsthorne@couchbase.com"
server = "172.23.104.149:19923"
```

Allocate cluster

```bash
> cbdyncluster allocate --num-nodes 3 --server-version 5.5.0
a004f847
> cbdyncluster ps
a004f847 [Owner: brett@couchbase.com, Creator: brett@couchbase.com, Timeout: 59m48s]
  dac14ff7eab9      node_1               5.5.0      172.23.111.210
  be8473e4b4d4      node_2               5.5.0      172.23.111.209
  72e31d4e0629      node_3               5.5.0      172.23.111.208
> cbdyncluster setup a004f847 --node kv --node kv --node kv --bucket default
http://172.23.111.210:8091
> cbdyncluster connstr a004f847
couchbase://172.23.111.210,172.23.111.209,172.23.111.208
> cbdyncluster rm a004f847
```

Pre release alias

```bash
--server-version 7.1-stable
```

GA release alias

```bash
--server-version 7.1-release
```

Capella Clusters

```bash
> cbdyncluster create-cloud --node kv,index,n1ql --node kv,index,n1ql --node kv,index,n1ql
e3f11b79
> cbdyncluster connstr e3f11b79 --ssl
couchbases://cb.fhdsf98hsd9fn.cloud.couchbase.com
```

AWS Clusters

```bash
> cbdyncluster allocate --num-nodes 3 --server-version 7.1-stable --platform ec2
43502742
> cbdyncluster ps
Clusters:
  43502742 [Owner: jake.rawsthorne@couchbase.com, Creator: jake.rawsthorne@couchbase.com, Timeout: 59m37s]
    i-0e969e29ba3a7c407  node_3               7.1.1-3070 ec2-3-144-127-195.us-east-2.compute.amazonaws.com dynclsr-43502742-node_3
    i-00a7751a64e1e03df  node_2               7.1.1-3070 ec2-18-220-187-47.us-east-2.compute.amazonaws.com dynclsr-43502742-node_2
    i-0bc686adef2f53133  node_1               7.1.1-3070 ec2-18-218-177-172.us-east-2.compute.amazonaws.com dynclsr-43502742-node_1
```

Use --help for more advanced usage
