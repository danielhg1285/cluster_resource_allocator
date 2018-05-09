// ResourceAllocator
// This program allows an optimal distribution of resources across a Pacemaker cluster taking as parameters
// the nodes RAM and CPU capacity and the ram and cpu used by each resource.

package main

import (
	"bufio"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/danyboy1104/tree"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
type ResourceAllocator struct {
	Tree
}

type Tree interface {
    Combine([]interface{},[]*Combination)
}*/

// A cluster node
type Node struct {
	name       string
	coresValue int
	ramValue   int
	dispar     bool
}

func NewNode(name string, ncores int, nram int, dispar bool) *Node {
	nod := new(Node)
	nod.name = name
	nod.coresValue = ncores
	nod.ramValue = nram
	nod.dispar = dispar
	// nod is initialized
	return nod
}

func (nod *Node) Name() string {
	return nod.name
}

func (nod *Node) CpuValue() int {
	return nod.coresValue
}

func (nod *Node) RamValue() int {
	return nod.ramValue
}

func (nod *Node) IsDispar() bool {
	return nod.dispar
}

// A combination node and its assigned resources
type Combination struct {
	listResources []interface{}
	totalCpuValue int
	totalRamValue int
	nodeAssigned  string
}

func NewCombination(rsc []interface{}, totalcores int, totalram int, nod string) *Combination {
	comb := new(Combination)
	comb.listResources = rsc
	comb.totalCpuValue = totalcores
	comb.totalRamValue = totalram
	comb.nodeAssigned = nod
	// comb is initialized
	return comb
}

// Gets resource list
func (comb *Combination) ListResources() []interface{} {
	return comb.listResources
}

func (comb *Combination) TotalCpuValue() int {
	return comb.totalCpuValue
}

func (comb *Combination) TotalRamValue() int {
	return comb.totalRamValue
}

func (comb *Combination) NodeAssigned() string {
	return comb.nodeAssigned
}

// A cluster resource
type Resource struct {
	name       string
	node       string
	coresValue int
	ramvalue   int
}

func NewResource(name string, node string, cores int, ramvalue int) interface{} {
	rsc := new(Resource)
	rsc.name = name
	rsc.node = node
	rsc.coresValue = cores
	rsc.ramvalue = ramvalue
	// rsc is initialized
	return rsc
}

func (rsc *Resource) Name() string {
	return rsc.name
}

func (rsc *Resource) SetNodeName(nodeName string) {
	rsc.node = nodeName
}

func (rsc *Resource) NodeName() string {
	return rsc.node
}

func (rsc *Resource) CpuValue() int {
	return rsc.coresValue
}

func (rsc *Resource) RamValue() int {
	return rsc.ramvalue
}

/* Returns a boolean which indicates if there exists a no assigned resource. */
func ExistUnasignedResource(listResources []interface{}) bool {
	for i := 0; i < len(listResources); i++ {
		if listResources[i].(*Resource).NodeName() == "" {
			return true
		}
	}
	return false
}

/* Returns a boolean which indicates if a resource is asigned */
func IsAsignedResource(rsc interface{}) bool {
	return rsc.(*Resource).node != ""
}

/* Returns a boolean which indicates if the combination is assigned. */
func (comb *Combination) IsCombinationAssigned() bool {
	return comb.NodeAssigned() != ""
}

/* Returns a boolean which indicates if a combination resource is already assigned. */
func (comb *Combination) IsResourceOfCombinationAssigned(listResources []interface{}) bool {
	for i := 0; i < len(comb.ListResources()); i++ {
		for j := 0; j < len(listResources); j++ {
			if comb.ListResources()[i].(*Resource).Name() == listResources[j].(*Resource).Name() {
				if IsAsignedResource(listResources[j]) {
					return true
				}
			}
		}
	}
	return false
}

/* Assign combination to node. */
func (comb *Combination) AssignCombinationToNode(nodeName string, listResources []interface{}) {
	comb.nodeAssigned = nodeName
	for i := 0; i < len(comb.ListResources()); i++ {
		for j := 0; j < len(listResources); j++ {
			if comb.ListResources()[i].(*Resource).Name() == listResources[j].(*Resource).Name() {
				listResources[j].(*Resource).SetNodeName(nodeName)
			}
		}
	}
}

/* Unassign combination to node. */
func (comb *Combination) UnAssignCombinationToNode(nodeName string, listResources []interface{}) {
	for i := 0; i < len(comb.ListResources()); i++ {
		for j := 0; j < len(listResources); j++ {
			combResource := comb.ListResources()[i]
			combResourceName := combResource.(*Resource).Name()
			resource := listResources[j]
			resourceName := resource.(*Resource).Name()
			if combResourceName == resourceName {
				listResources[j].(*Resource).SetNodeName("")
			}
		}
	}
	comb.nodeAssigned = ""
}

// Prints resource on nodes
func PrintResourceMatrix(list []interface{}) {
	for i := 0; i < len(list); i++ {
		resources := list[i].(*Resource)
		fmt.Println(resources.Name())
	}
}

// Prints resources
func PrintResources(listado []interface{}) {
	for i := 0; i < len(listado); i++ {
		fmt.Println(listado[i].(*Resource).Name(), listado[i].(*Resource).NodeName())
	}
}

// Prints node resource combinations
func PrintCombinations(listado []*Combination) {
	for i := 0; i < len(listado); i++ {
		cadena := ""
		for j := 0; j < len(listado[i].ListResources()); j++ {
			cadena += listado[i].ListResources()[j].(*Resource).Name()
		}
		fmt.Println(cadena)
	}
}

// This function create a Matrix of all possible combinations of resource per nodes
func Combine(treeNode *tree.Tree, resources []interface{}, combinations *[]*Combination) {
	for i := 0; i < len(resources); i++ {
		if treeNode.IsRoot() {
			var slice []interface{}
			slice = append(slice, resources[i])
			treeNode.SetChildren(i, tree.NewTree(slice, false, resources[i].(*Resource).CpuValue(), resources[i].(*Resource).RamValue(), ""))
			comb := NewCombination(slice, resources[i].(*Resource).CpuValue(), resources[i].(*Resource).RamValue(), "")
			*combinations = append(*combinations, comb)

		} else {
			childData := make([]interface{}, len(treeNode.Data()))
			copy(childData, treeNode.Data())
			childData = append(childData, resources[i])
			sumaCpu := 0
			sumaRam := 0
			for j := 0; j < len(childData); j++ {
				sumaCpu += childData[j].(*Resource).CpuValue()
				sumaRam += childData[j].(*Resource).RamValue()
			}
			treeNode.SetChildren(i, tree.NewTree(childData, false, sumaCpu, sumaRam, ""))
			comb := NewCombination(childData, sumaCpu, sumaRam, "")
			*combinations = append(*combinations, comb)
			if len(resources) < 2 {
				return
			}
		}
		Combine(treeNode.Children()[i], resources[i+1:len(resources)], combinations)
	}
}

// This function find a valid solution in all possible combination if exists
// based on the values of cpu and ram for each resource. It use Backtracking
func DistributeResources(listNodes []*Node, listResources []interface{}, listCombinations []*Combination) bool {
	// If there is no unassigned resource, we are done
	if !ExistUnasignedResource(listResources) {
		return true // success!
	}

	for j := 0; j < len(listCombinations); j++ {
		if listCombinations[j].TotalCpuValue() <= listNodes[0].CpuValue() && listCombinations[j].TotalRamValue() <= listNodes[0].RamValue() {
			if !listCombinations[j].IsCombinationAssigned() {
				if !listCombinations[j].IsResourceOfCombinationAssigned(listResources) {
					listCombinations[j].AssignCombinationToNode(listNodes[0].Name(), listResources)
					if len(listNodes) > 1 {
						if DistributeResources(listNodes[1:len(listNodes)], listResources, listCombinations) {
							return true
						}
					} else {
						if !ExistUnasignedResource(listResources) {
							return true // success!
						}
					}
					// failure, unmake & try again
					if listCombinations[j].IsCombinationAssigned() {
						listCombinations[j].UnAssignCombinationToNode(listNodes[0].Name(), listResources)
					}
				}
			}
		}
	}
	return false // this triggers backtracking
}

// Return true if the resource is on the matrix
func ExistResourceNodeOnMatrixByNode(resource interface{}, matrixByNode [][]interface{}) bool {
	for i := range matrixByNode {
		if matrixByNode[i][0].(*Resource).NodeName() == resource.(*Resource).NodeName() {
			return true
		}
	}
	return false
}

// Assign a resource to a Node
func InsertResourceOnMatrixByNode(resource interface{}, matrixByNode [][]interface{}) {
	for i := range matrixByNode {
		if matrixByNode[i][0].(*Resource).NodeName() == resource.(*Resource).NodeName() {
			matrixByNode[i] = append(matrixByNode[i], resource)
		}
	}
}

// Create a matrix of only resources
func CreateResourceMatrixByNode(listResources []interface{}) [][]interface{} {
	var matrixByNode [][]interface{}
	for i := 0; i < len(listResources); i++ {
		if !ExistResourceNodeOnMatrixByNode(listResources[i], matrixByNode) {
			var slice []interface{}
			slice = append(slice, listResources[i])
			matrixByNode = append(matrixByNode, slice)
		} else {
			InsertResourceOnMatrixByNode(listResources[i], matrixByNode)
		}
	}
	return matrixByNode
}

func PrintResourceNodeMatrix(resourceMatrix [][]interface{}) {
	for i := range resourceMatrix {
		for j := range resourceMatrix[i] {
			fmt.Println(resourceMatrix[i][j].(*Resource).Name(), resourceMatrix[i][j].(*Resource).NodeName())
		}
	}
}

func IsNodeOfResourceDispar(listNodes []*Node, resource interface{}) bool {
	for i := range listNodes {
		if listNodes[i].Name() == resource.(*Resource).NodeName() {
			if listNodes[i].IsDispar() {
				return true
			}
		}
	}
	return false
}

// Create the cluster rules. Based on Pacemaker colocation rules.
func CreateClusterRules(listNodes []*Node, resourceMatrix [][]interface{}) []string {
	const score int = 10000
	var rules []string
	for i := range resourceMatrix {
		// If resource of resourceMatrix on i is dispar add rule of node resource location
		if IsNodeOfResourceDispar(listNodes, resourceMatrix[i][0]) {
			targetResource := resourceMatrix[i][0].(*Resource).Name()
			targetNode := resourceMatrix[i][0].(*Resource).NodeName()
			crmLocationCommand := "cluster_prefers_node --node " + targetNode + " --resource " + targetResource + " --score " + strconv.Itoa(score)
			rules = append(rules, crmLocationCommand)
		}
		for j := i + 1; j < len(resourceMatrix); j++ {
			for k := range resourceMatrix[j] {
				// if i equal j go to next row
				if j == i {
					break
				} else {
					// Else add rule resourceMatrix on [i][0] anticolocacion resourceMatrix on [j][k]
					firstResource := resourceMatrix[i][0].(*Resource).Name()
					secondResource := resourceMatrix[j][k].(*Resource).Name()
					crmColocationCommand := "cluster_anticolocate_resources --first-resource " + firstResource + " --second-resource " + secondResource + " --score " + strconv.Itoa(-1*score)
					rules = append(rules, crmColocationCommand)
				}
			}
		}

	}
	return rules
}

func PrintRules(rules []string) {
	for i := range rules {
		fmt.Println(rules[i])
	}
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Execute command that get the cluster node names and resources
	cmdName := dir + "/cluster_get_nodes_resources.sh"
	cmd := exec.Command(cmdName)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	var nodes []string
	var resources []string
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			lineFields := strings.Fields(scanner.Text())
			if lineFields[0] == "node" {
				nodes = append(nodes, lineFields[1])
			} else {
				resources = append(resources, lineFields[1])
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	var listNodes []*Node
	// Here we assign the value for ram and cpu for every server
	for i := range nodes {
		nodeCores := (&prompter.Prompter{
			Message: "Enter numer of cores for node " + nodes[i],
			Default: "10",
		}).Prompt()
		nodeRam := (&prompter.Prompter{
			Message: "Enter Ram in GB for node " + nodes[i],
			Default: "10",
		}).Prompt()
		dispar := prompter.YN("Is the node "+nodes[i]+" dispar", true)
		valNodeCores, _ := strconv.Atoi(nodeCores)
		valNodeRam, _ := strconv.Atoi(nodeRam)
		listNodes = append(listNodes, NewNode(nodes[i], valNodeCores, valNodeRam, dispar))
	}

	var listResources []interface{}
	// Here we request to user the value for ram and cpu for every resource
	for i := range resources {
		rscCores := (&prompter.Prompter{
			Message: "Enter numer of cores for resource " + resources[i],
			Default: "0",
		}).Prompt()
		rscRam := (&prompter.Prompter{
			Message: "Enter used RAM in GB for resource " + resources[i],
			Default: "0",
		}).Prompt()
		valRscCores, _ := strconv.Atoi(rscCores)
		valRscRam, _ := strconv.Atoi(rscRam)
		listResources = append(listResources, NewResource(resources[i], "", valRscCores, valRscRam))
	}

	//listNodes := []*Node{{"node1", 12, 8, true}, {"node2", 20, 9, false}, {"node3", 10, 20, true}}
	//listResources := []*Resource{{"rsc1", "", 4, 3}, {"rsc2", "", 10, 20}, {"rsc3", "", 6, 3}, {"rsc4", "", 7, 4}, {"rsc5", "", 5, 3}, {"rsc6", "", 8, 2}, {"rsc7", "", 2, 2}}

	var listCombinations []*Combination
	var genericResourceList []interface{}
	genericResourceList = make([]interface{}, len(listResources))
	for i, d := range listResources {
		genericResourceList[i] = d
	}

	treeNode := tree.NewTree(genericResourceList, true, 0, 0, "")
	Combine(treeNode, genericResourceList, &listCombinations)

	if DistributeResources(listNodes, genericResourceList, listCombinations) {
		//PrintResources(genericResourceList)
		resourceMatrixByNode := CreateResourceMatrixByNode(genericResourceList)
		PrintResourceNodeMatrix(resourceMatrixByNode)
		// Here the cluster rules are created to be executed later
		//clusterRules := CreateClusterRules(listNodes, resourceMatrixByNode)
		//PrintRules(clusterRules)
	} else {
		fmt.Println("Ther's no capacity to allocate resources on existing nodes")
	}

	//PrintCombinations(listCombinations)
	//var matrix [][]interface{}
	//treeNode.Preorder(&matrix)

}
