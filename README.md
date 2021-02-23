# Biosplitter

A reference implementation for the proposed splitting logic of the [KubeIT](https://github.com/KubeITerator/KubeIT) project. This Splitter splits FASTA files by either number of records or byte-size.

## Behavior

Biosplitter uses the environment variables

- `DATASOURCE`: A URL to a file that should be distributed
- `PARAMS`:  A JSON object that describes parameters that specify the actual splitting

to determine suitable locations for the splitting of the input file. These locations are outputted as a JSON formatted list
that looks like this:

```json

  [
    { "index": 0, "range": "Range:bytes=0-1000" },
    { "index": 1, "range": "Range:bytes=1001-2000" }
  ]

```

The `index` value is an incrementing number. `range` refers to the string `Range:bytes=START-STOP` that specifier a HTML Range HEADER in curl.

This behaviour is distributed in a [docker container](https://hub.docker.com/repository/docker/stanni/biosplitter)

#### PARAMS

Biosplitter currently accepts two params to determine suitable splitting positions.
`maxrecords` and `bytesize`. `maxrecords` defines an upper limit for the number of records, while `bytesize` defines a approximate size in bytes per chunk.
The PARAM envvar must be JSON formatted, example:

```json
{
  "maxrecord": 1,
  "bytesize": 100000
}
```

Only one param must be specified, if both are specified biosplitter prefers to split for the `maxrecord` factor.

### SplitterInterface


To create your own Splitting logic you must recreate the above behaviour for your own container. If you use Go for your container
you can use the [SplitterInterface](/logic/SplitterInterface.go) as interface for your operation. 

