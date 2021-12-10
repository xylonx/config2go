# config2go

a template repo to help me start a golang project quickly. 

Generally, use yaml config file with paring by viper. Use cobra to make the project start as a normal command line tool.

## how to install

`go install github.com/xylonx/config2go@latest`

## How to use

`config2go -s ${sourceConfig} -t ${targetConfig} -p ${package} -t ${tags}`

for more information, run `config2go -h`.