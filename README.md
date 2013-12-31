# nethogs-parser

A parser to summarize the nethogs trace mode output

# Nethogs output parser

Nethogs is a `net top` tool. [http://nethogs.sourceforge.net/]

# Usage - go script

+ `go run hogs.go <options> [file]...`
+ `-type=csv|pretty`
+ `-datatable=<datatable_location>`
+ `-class=<datatable_classname>`
+ `-cpuprofile=<filename>.prof`

+ If you'd like to deploy it somewhere, `go build hogs.go` and `./hogs <options>`

# Examples

## Already existing output - Go [ master ]

+ `./hogs -type=csv output1 output2 output3`
+ `./hogs -datatable localhost:4200 -class nethogsbw output1 output2`
+ `./hogs -type=pretty -datatable localhost:4200 -class mynethogs output1`
