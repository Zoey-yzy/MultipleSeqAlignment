package main

import (
	
	"fmt"
	"os"
)

//delete bc didn't work

func WriteNewickToFile(t Tree, fileDest string, fileName string) {

	newickString := ToNewick(t)

	F, err := os.Create(fileDest + "/" + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = F.WriteString(newickString)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = F.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func ToNewick(tree Tree) string {
	return "(" + subtreeNewickAges(tree[len(tree)-1]) + ");"
}

func subtreeNewickAges(node *Node) string {
	fmt.Println("node.label:",node.label,len(node.neighbors))
	if len(node.neighbors) == 0 {
		return ""

	}else if node.neighbors[0]== nil {
		return node.label + ":" + fmt.Sprintf("%.2f", node.distance)
	} else {
		return "(" + subtreeNewickAges(node.neighbors[0]) + "," + subtreeNewickAges(node.neighbors[0]) + "):" + fmt.Sprintf("%.2f", node.distance)
	}
}
