# nethogs-parser

A parser to summarize the nethogs trace mode output

# Nethogs output parser

Nethogs is a `net top` tool. [http://nethogs.sourceforge.net/]

# Usage - go script

+ `go run hogs.go [-type=csv|pretty] [-datatable=<hostname>:<port>] [file]...`
+ If you'd like to deploy it somewhere, `go build hogs.go` and `./hogs <options>`

# Usage - python script [ Will be DEPRECATED ]

`python nethogs.bw.py <timeout_in_seconds>`

# Examples

## Already existing output - Go [ master ]

+ `./hogs -type=csv output1 output2 output3`
+ `./hogs -datatable localhost:4200 output1 output2`
+ `./hogs -type=pretty -datatable localhost:4200 output1`

## Already existing output - Python [ Will be DEPRECATED ]

+ `nethogs -t eth1 > nethogs.out`
+ `python nethogs.bw.py 1000 < nethogs.out`
+ # 1000 is just a timeout

## Listen on nethogs output for some <timeout>  - Python [ Will be DEPRECATED ]

+ `nethogs -t eth1 | python nethogs.bw.py 3600`
+ # Nethogs listens on eth1 and pipes output to nethogs.bw.py which parses them and aggregates the data
+ # After `3600s`, the python program exits sending quit signal to nethogs too, and prints the output
