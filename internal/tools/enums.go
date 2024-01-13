package tools

// SA1029: should not use built-in type string as key for value; define your own type to avoid collisions (staticcheck).
type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	bufferSize     = 2048
	SudoPassCtxKey = contextKey("sudoPassword")

	echoCmd     = "echo"
	sudoCmd     = "sudo"
	topCmd      = "top"
	iostatCmd   = "iostat"
	grepCmd     = "grep"
	diskFreeCmd = "df"
	netstatCmd  = "netstat"
	ssCmd       = "ss"

	percentPattern  = `.*?(\d{1,3}),(\d{1,3}).*?`
	diskFreePattern = `.*?\s(\w+)\s.*?(\d+)%\s.*`
	netstatPattern  = `^(\w+)\s+.*:(\d+)\s+.*\s+(\d+)\s+\d+\s+(\d+)\/(.*)$`
	ssPattern       = `^([\w\-]+)\s+.*`
)

var (
	topArgs             = []string{"-b", "-n1"}
	iostatArgs          = []string{"-d", "-k"}
	diskSizeArgs        = []string{"-mT"}
	diskInodesArgs      = []string{"-iT"}
	netstatArgs         = []string{"-lntupe"}
	netstatWithSudoArgs = []string{"-S", netstatCmd, "-lntupe"}
	ssTCPArgs           = []string{"-ta"}
	ssUDPArgs           = []string{"-ua"}

	loadAverageArgs = []string{"average"}
	cpuLoadArgs     = []string{"Cpu"}
	diskLoadArgs    = []string{"sda"}
)