package types

type Enum interface {
	String() string
	FromI(i int) Enum
	IsValid() bool
}

type Format int

const InvalidString = "invalid"

const (
	SVG Format = iota
	PNG
)

func (f Format) String() string {
	switch f {
	case SVG:
		return "svg"
	case PNG:
		return "png"
	default:
		return InvalidString
	}
}

func (f Format) FromI(i int) Enum {
	return Format(i)
}

func (f Format) IsValid() bool {
	return f.String() != InvalidString
}

func All[T Enum]() []string {
	options := []string{}

	t := *new(T)
	i := 0

	for {
		thisOpt := t.FromI(i)
		if thisOpt.String() != InvalidString {
			options = append(options, thisOpt.String())
			i++
		} else {
			break
		}
	}

	return options
}

type LineEnding int

const (
	LF LineEnding = iota
	CRLF
)

func (l LineEnding) String() string {
	switch l {
	case LF:
		return "LF"
	case CRLF:
		return "CRLF"
	default:
		return InvalidString
	}
}

func (l LineEnding) FromI(i int) Enum {
	return LineEnding(i)
}

func (l LineEnding) IsValid() bool {
	return l.String() != InvalidString
}

type AppConfig struct {
	Input     string
	Output    string
	Format    Format
	Linewidth int
	Ending    LineEnding
	Verbose   bool
}

func GetFormat(i string) Format {
	switch i {
	case "svg", "SVG":
		return SVG
	case "png", "PNG":
		return PNG
	default:
		return Format(-1)
	}
}

func GetEnding(i string) LineEnding {
	switch i {
	case "crlf", "CRLF":
		return CRLF
	case "lf", "LF":
		return LF
	default:
		return LineEnding(-1)
	}
}
