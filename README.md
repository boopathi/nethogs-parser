# Nethogs parser

A parser to summarize the nethogs trace mode output

[![Build Status](https://travis-ci.org/boopathi/nethogs-parser.png?branch=master)](https://travis-ci.org/boopathi/nethogs-parser)

# Nethogs

Nethogs is a `net top` tool. [http://nethogs.sourceforge.net/]

# Make

+ `go build -o hogs hogs.go`

# Usage - go script

+ `./hogs <options> [file]...`
+ `-type=csv|pretty`
+ `-datatable=<datatable_location>`
+ `-class=<datatable_classname>`
+ `-cpuprofile=<filename>.prof`

# Examples

+ `./hogs -type=csv output1 output2 output3`
+ `./hogs -datatable localhost:4200 -class nethogsbw output1 output2`
+ `./hogs -type=pretty -datatable localhost:4200 -class mynethogs output1`
