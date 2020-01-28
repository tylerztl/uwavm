package gas

const (
	maxResourceLimit = 0xFFFFFFFF
	CpuRate          = 1000
	MemRate          = 1000000
	DiskRate         = 1
	FeeRate          = 1
)

// Limits describes the usage or limit of resources
type Limits struct {
	Cpu    int64
	Memory int64
	Disk   int64
	Fee    int64
}

// TotalGas converts resource to gas
func (l *Limits) TotalGas() int64 {
	cpuGas := roundup(l.Cpu, CpuRate)
	memGas := roundup(l.Memory, MemRate)
	diskGas := roundup(l.Disk, DiskRate)
	feeGas := roundup(l.Fee, FeeRate)
	return cpuGas + memGas + diskGas + feeGas
}

// Add accumulates resource limits, returns self.
func (l *Limits) Add(l1 Limits) *Limits {
	l.Cpu += l1.Cpu
	l.Memory += l1.Memory
	l.Disk += l1.Disk
	l.Fee += l1.Fee
	return l
}

// Sub sub limits from l
func (l *Limits) Sub(l1 Limits) *Limits {
	l.Cpu -= l1.Cpu
	l.Memory -= l1.Memory
	l.Disk -= l1.Disk
	l.Fee -= l1.Fee
	return l
}

// Exceed judge whether resource exceeds l1
func (l Limits) Exceed(l1 Limits) bool {
	return l.Cpu > l1.Cpu ||
		l.Memory > l1.Memory ||
		l.Disk > l1.Disk ||
		l.Fee > l1.Fee
}

// MaxLimits describes the maximum limit of resources
var MaxLimits = Limits{
	Cpu:    maxResourceLimit,
	Memory: maxResourceLimit,
	Disk:   maxResourceLimit,
	Fee:    maxResourceLimit,
}

func roundup(n, scale int64) int64 {
	if scale == 0 {
		return 0
	}
	return (n + scale - 1) / scale
}
