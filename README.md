This is a CLI tool for interfacing with a running cbdynclusterd.

Example Usage (as of Aug 22, 2018):
```bash
> cbdyncluster allocate --num-nodes 3 --server-version 5.5.0
a004f847
> cbdyncluster ps
a004f847 [Owner: brett@couchbase.com, Creator: brett@couchbase.com, Timeout: 59m48s]
  dac14ff7eab9      node_1               5.5.0      172.23.111.210
  be8473e4b4d4      node_2               5.5.0      172.23.111.209
  72e31d4e0629      node_3               5.5.0      172.23.111.208
> cbdyncluster connstr a004f847
couchbase://172.23.111.210,172.23.111.209,172.23.111.208
> cbdyncluster rm a004f847
```