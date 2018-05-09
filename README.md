# cluster_resource_allocator
A tool for automatic distribution of cluster resources accross all Pacemaker nodes in an optimal way, taking into account the total cluster hardware and the hardware consumed by resources, making easy the automatic configuration when there's need to scale up.

##Requirements
### golang installed

### Go packages
bufio
fmt
os
github.com/Songmu/prompter
github.com/danyboy1104/tree
filepath


##Installing
###1. Build cluster_resource_allocator.go using golang
###2. Copy script cluster_get_nodes_resources.sh on same directory as cluster_resource_allocator binary.

##Usage
###3. Execute cluster_resource_allocator binary
```bash
./cluster_resource_allocator
```
###4. Enter each node value for CPU and RAM
###5. Enter each resource value for CPU and RAM




