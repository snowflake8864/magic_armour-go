package tree

/**
var nodes []*tree.Node
for _, model := range models {
	nodes = append(nodes, &tree.Node{
		ID:       model.IdSingle,
		ParentId: model.ParentId,
		Data: map[string]interface{}{
			"name":       model.Name,
			"code":       model.Code,
			"sort":       model.Sort,
			"is_show":    model.IsShow,
			"is_root":    model.IsRoot,
			"remark":     model.Remark,
			"api":        model.Api,
			"level":      model.Level,
			"created_at": model.CreatedAt,
		},
		Child: []*tree.Node{},
	})
}
tree.RealTree(nodes)
*/

type Node struct {
	ID       int64                  `json:"id"`
	ParentId int64                  `json:"parentId"`
	Data     map[string]interface{} `json:"data"`
	Child    []*Node                `json:"children"`
}

func MakeTree(nodes []*Node, node *Node) {
	childNodes := findChilds(node, nodes)
	for _, childNode := range childNodes {
		node.Child = append(node.Child, childNode)
		if hasChild(childNode, nodes) {
			MakeTree(nodes, childNode)
		}
	}
}

func findChilds(node *Node, nodes []*Node) (ret []*Node) {
	for _, currNode := range nodes {
		if currNode.ParentId == int64(node.ID) {
			ret = append(ret, currNode)
		}
	}
	return
}

func hasChild(node *Node, nodes []*Node) bool {
	for _, currNode := range nodes {
		if currNode.ParentId == int64(node.ID) {
			return true
		}
	}
	return false
}

func RealTree(nodes []*Node) []*Node {
	var firstNode = Node{
		ID:       0,
		Data:     nil,
		ParentId: -1,
		Child:    []*Node{},
	}
	if len(nodes) > 0 {
		MakeTree(nodes, &firstNode)
	}
	return firstNode.Child
}
