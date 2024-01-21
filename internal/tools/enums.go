package tools

// SA1029: should not use built-in type string as key for value; define your own type to avoid collisions (staticcheck).
type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	bufferSize     = 2048
	zeroNumber     = 0.0
	SudoPassCtxKey = contextKey("sudoPassword")

	echoCmd     = "echo"
	sudoCmd     = "sudo"
	topCmd      = "top"
	iostatCmd   = "iostat"
	grepCmd     = "grep"
	diskFreeCmd = "df"
	netstatCmd  = "netstat"
	ssCmd       = "ss"

	percentWithComaPattern   = `.*?(\d+),(\d+).*?`
	percentWithPointPattern  = `.*?(\d+)\.(\d+).*?`
	loadAveragePatternFormat = ".*average:%s,.*"
	cpuLoadPatternFormat     = "%suser,%ssys,%sidle.*"
	diskLoadPatternFormat    = ".*sda%s%s%s"
	diskFreePattern          = `.*?\s(\w+)\s.*?(\d+)%\s.*`
	netstatPattern           = `^(\w+)\s+.*:(\d+)\s+.*\s+(\d+)\s+\d+\s+(\d+)\/(.*)$`
	ssPattern                = `^([\w\-]+)\s+.*`
)

var (
	topArgs             = []string{"-n1"}
	iostatArgs          = []string{"-d", "-k"}
	diskSizeArgs        = []string{"-m"}
	diskInodesArgs      = []string{"-i"}
	netstatArgs         = []string{"-lntupe"}
	netstatWithSudoArgs = []string{"-S", netstatCmd, "-lntupe"}
	ssTCPArgs           = []string{"-ta"}
	ssUDPArgs           = []string{"-ua"}

	loadAverageArgs = []string{"'Load Avg'"}
	cpuLoadArgs     = []string{"'CPU usage'"}
	diskLoadArgs    = []string{"sda"}
)
