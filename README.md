# aleo_exporter

aleo_exporter exports basic monitoring data from a Aleo node.

## Metrics
- **latest_height** - the latest block height
- **peer_count** - the number of peers connected to the node

## Command line arguments

You typically only need to set the RPC URL, pointing to one of your own nodes:

    ./aleo_exporter -rpcURI=http://yournode:3033

```
Usage of ./aleo_exporter:
  -add_dir_header
    	If true, adds the file directory to the header of the log messages
  -addr string
    	Listen address (default ":8080")
  -alsologtostderr
    	log to standard error as well as files
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -log_file string
    	If non-empty, use this log file
  -log_file_max_size uint
    	Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
  -logtostderr
    	log to standard error instead of files (default true)
  -one_output
    	If true, only write logs to their native severity level (vs also writing to each lower severity level
  -rpcURI string
    	Aleo RPC URI
  -skip_headers
    	If true, avoid header prefixes in the log messages
  -skip_log_headers
    	If true, avoid headers when opening log files
  -stderrthreshold value
    	logs at or above this threshold go to stderr (default 2)
  -v value
    	number for the log level verbosity
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```
