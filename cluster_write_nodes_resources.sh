#!/bin/bash
# cluster_write_nodes_resources
#

# The functions to obtain the cluster nodes and resources are commented for simplicity 
: <<'end_long_comment'
# List cluster nodes
online_hosts_list | tr ' ' '\n' | sort | \
while read host ; do
  echo "node $host"
done
# List cluster resources
movable_rscs | sort | \
while read resource; do
  echo "resource $resource"
done
end_long_comment

# Here we use an example output of nodes and resources
cat << EOT
node node1
node node2
node node3
resource rsc1
resource rsc2
resource rsc3
resource rsc4
resource rsc5
resource rsc6
resource rsc7
EOT
