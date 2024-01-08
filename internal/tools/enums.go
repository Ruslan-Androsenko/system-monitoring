package tools

const (
	topCmd      = "top"
	iostatCmd   = "iostat"
	grepCmd     = "grep"
	diskFreeCmd = "df"

	percentPattern  = `.*?(\d{1,3}),(\d{1,3}).*?`
	diskFreePattern = `.*?\s(\w+)\s.*?(\d+)%\s.*`
)

var (
	topArgs        = []string{"-b", "-n1"}
	iostatArgs     = []string{"-d", "-k"}
	diskSizeArgs   = []string{"-mT"}
	diskInodesArgs = []string{"-iT"}

	loadAverageArgs = []string{"average"}
	cpuLoadArgs     = []string{"Cpu"}
	diskLoadArgs    = []string{"sda"}
)
