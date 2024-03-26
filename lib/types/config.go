package types

type AppConfig struct {
	Input        string
	Output       string
	Format       Format
	Linewidth    int
	Ending       LineEnding
	Verbose      bool
	XRes         int
	YRes         int
	Simres       int
	Algorithm    Algorithm
	RefineCycles int
}
