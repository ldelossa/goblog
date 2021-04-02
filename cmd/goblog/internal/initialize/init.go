package initialize

import (
	"context"
	"flag"
	"math"

	"github.com/ldelossa/goblog/pkg/golog"
)

var initFS = flag.NewFlagSet("init", flag.ExitOnError)

var initFlags = struct{}{}

func Initialize(ctx context.Context) {
	var root = buildDTree()
	err := root.Execute(ctx)
	if err != nil {
		golog.Fatal("%v\n", err)
	}
}

// Initialization is driven off whether a goblog home directory
// can be found.
//
// If the home directory cannot be found the user will be asked to
// define one, provide a remote and branch to clone the goblog source
// into this home directory, and build a new goblog binary.
//
// If a home directory is found we check to see if a git repo
// exists in it.
//
//			 				     Decision Tree for Initialization
//H1		 						 [  0 HomeExistsDecision  ]
//			 					   no/                        \yes
//H2			 	      [1 GitCloneDecision]                  [2 CheckGitRepoDecision]
//  			 	      /         \                                /                 \
//H3			        [3]     [4 Synchronize]           [5 GitCloneDecision]          [6 Synchronize]
//  			      /     \         /    \                    /     \                        /     \
//H4			    [7]     [8]     [9] [10 BuildDecision]   [11] [12 Synchronize]           [13]  [14 UpgradeDecision]
//                 /   \    / \     / \        /    \        /  \       /  \                 /  \         /  \
//H5	       [15]  [16][17][18][19][20]     [21] [22]   [23][24]   [25][26 BuildDecision] [27][28]     [29] [30]

func buildDTree() *Decision {
	height := 5.0
	N := math.Exp2(height) - 1
	nodes := make([]*Decision, int(N), int(N))
	nodes[0] = &Decision{Exec: HomeExists}
	nodes[1] = &Decision{Exec: GitClone}
	nodes[2] = &Decision{Exec: CheckGitRepo}
	nodes[4] = &Decision{Exec: Synchronize}
	nodes[5] = &Decision{Exec: GitClone}
	nodes[6] = &Decision{Exec: Synchronize}
	nodes[10] = &Decision{Exec: Build}
	nodes[12] = &Decision{Exec: Synchronize}
	nodes[14] = &Decision{Exec: Upgrade}
	nodes[26] = &Decision{Exec: Build}

	// iteratate and link children
	for i, node := range nodes {
		if (2*i)+1 > len(nodes)-1 {
			break
		}
		if node == nil {
			continue
		}
		node.AddNo(
			nodes[(2*i)+1],
		)
		node.AddYes(
			nodes[(2*i)+2],
		)
	}
	return nodes[0]
}
