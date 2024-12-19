package port_whether_occupy

func CreateFactory() *Curd {
	return &Curd{
		bitArr:    make([]int64, 65536/32),
		separator: 32,
	}
}

type Curd struct {
	bitArr    []int64
	separator int64
}

// 端口所在索引
func (crd *Curd) toIndex(port int64) int64 {
	return port / crd.separator
}

// 端口所在位
func (crd *Curd) toBit(port int64) int64 {
	return port % crd.separator
}

// SetPort 将位图中的特定位设置为1，表示该端口已使用
func (crd *Curd) SetPort(port int64) {
	index := crd.toIndex(port)
	bit := crd.toBit(port)

	crd.bitArr[index] = crd.bitArr[index] | 1<<bit
}

// CheckPort 将位检查端口是否在位图中被标记为已使用, 未使用返回true
func (crd *Curd) CheckPort(port int64) bool {
	index := crd.toIndex(port)
	bit := crd.toBit(port)
	return (crd.bitArr[index] & (1 << bit)) == 0
}
