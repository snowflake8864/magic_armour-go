package port_whether_occupy

func CreateFactory() *Curd {
	return &Curd{
		bitArr:    make([]int, 63356/32),
		separator: 32,
	}
}

type Curd struct {
	bitArr    []int
	separator int
}

// 端口所在索引
func (crd *Curd) toIndex(port int) int {
	return port / crd.separator
}

// 端口所在位
func (crd *Curd) toBit(port int) int {
	return port % crd.separator
}

// SetPort 将位图中的特定位设置为1，表示该端口已使用
func (crd *Curd) SetPort(port int) {
	index := crd.toIndex(port)
	bit := crd.toBit(port)

	crd.bitArr[index] = crd.bitArr[index] | 1<<bit
}

// CheckPort 将位检查端口是否在位图中被标记为已使用, 未使用返回true
func (crd *Curd) CheckPort(port int) bool {
	index := crd.toIndex(port)
	bit := crd.toBit(port)

	return (crd.bitArr[index] & (1 << bit)) == 0
}
