# nethogs-parser

A parser to summarize the nethogs trace mode output

# Nethogs output parser

Nethogs is a `net top` tool. [http://nethogs.sourceforge.net/]

# Usage

`python nethogs.bw.py <timeout_in_seconds>`

# Examples

## Already existing output

+ `nethogs -t eth1 > nethogs.out`
+ `python nethogs.bw.py 1000 < nethogs.out`
+ # 1000 is just a timeout

## Listen on nethogs output for some <timeout> 

+ `nethogs -t eth1 | python nethogs.bw.py 3600`
+ # Nethogs listens on eth1 and pipes output to nethogs.bw.py which parses them and aggregates the data
+ # After `3600s`, the python program exits sending quit signal to nethogs too, and prints the output
