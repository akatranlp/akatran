package bytesize

import (
	"fmt"
	"regexp"
	"strings"
)

type ByteSize uint64

const (
	B ByteSize = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB

	KB ByteSize = 1000
	MB ByteSize = KB * 1000
	GB ByteSize = MB * 1000
	TB ByteSize = GB * 1000
	PB ByteSize = TB * 1000
	EB ByteSize = PB * 1000
)

var mapUnitToByteSize = map[string]ByteSize{
	"B":   B,
	"KB":  KB,
	"MB":  MB,
	"GB":  GB,
	"TB":  TB,
	"PB":  PB,
	"EB":  EB,
	"KiB": KiB,
	"MiB": MiB,
	"GiB": GiB,
	"TiB": TiB,
	"PiB": PiB,
	"EiB": EiB,
	"":    B,
}

var mapByteSizeToUnit = map[ByteSize]string{
	B:   "B",
	KB:  "KB",
	MB:  "MB",
	GB:  "GB",
	TB:  "TB",
	PB:  "PB",
	EB:  "EB",
	KiB: "KiB",
	MiB: "MiB",
	GiB: "GiB",
	TiB: "TiB",
	PiB: "PiB",
	EiB: "EiB",
}

func (b ByteSize) Format(format string, unit string) string {
	return b.format(format, unit)
}

func (b ByteSize) format(format string, unit string) string {
	var unitSize ByteSize
	if unit != "" {
		var ok bool
		unitSize, ok = mapUnitToByteSize[unit]
		if !ok {
			panic(fmt.Sprintf("invalid unit: %s", unit))
		}

		paramMatcher := regexp.MustCompile("%.*f")
		if !paramMatcher.MatchString(format) {
			panic(fmt.Sprintf("invalid format for specific unit: %s", format))
		}
		return fmt.Sprintf(format, float64(b)/float64(unitSize)) + mapByteSizeToUnit[unitSize]

	}
	switch {
	case b >= EiB && b%EiB == 0:
		return fmt.Sprintf("%dEiB", b/EiB)
	case b >= EB && b%EB == 0:
		return fmt.Sprintf("%dEB", b/EB)
	case b >= PiB && b%PiB == 0:
		return fmt.Sprintf("%dPiB", b/PiB)
	case b >= PB && b%PB == 0:
		return fmt.Sprintf("%dPB", b/PB)
	case b >= TiB && b%TiB == 0:
		return fmt.Sprintf("%dTiB", b/TiB)
	case b >= TB && b%TB == 0:
		return fmt.Sprintf("%dTB", b/TB)
	case b >= GiB && b%GiB == 0:
		return fmt.Sprintf("%dGiB", b/GiB)
	case b >= GB && b%GB == 0:
		return fmt.Sprintf("%dGB", b/GB)
	case b >= MiB && b%MiB == 0:
		return fmt.Sprintf("%dMiB", b/MiB)
	case b >= MB && b%MB == 0:
		return fmt.Sprintf("%dMB", b/MB)
	case b >= KiB && b%KiB == 0:
		return fmt.Sprintf("%dKiB", b/KiB)
	case b >= KB && b%KB == 0:
		return fmt.Sprintf("%dKB", b/KB)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func (b ByteSize) String() string {
	return b.format("", "")
}

func Parse(s string) (ByteSize, error) {
	var num uint64
	var unit string

	if _, err := fmt.Sscanf(s, "%d%s", &num, &unit); err != nil {
		return 0, err
	}

	mult, ok := mapUnitToByteSize[strings.TrimSpace(unit)]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %s", unit)
	}
	return ByteSize(num) * mult, nil
}

// Satisfy the flag package  Value interface.
func (b *ByteSize) Set(s string) error {
	bs, err := Parse(s)
	if err != nil {
		return err
	}
	*b = bs
	return nil
}

// Satisfy the pflag package Value interface.
func (b *ByteSize) Type() string { return "byte_size" }

// Satisfy the encoding.TextUnmarshaler interface.
func (b *ByteSize) UnmarshalText(text []byte) error {
	return b.Set(string(text))
}

// Satisfy the flag package Getter interface.
func (b *ByteSize) Get() interface{} { return ByteSize(*b) }

func (b ByteSize) FromBytes() uint64 { return uint64(b) }
func (b ByteSize) FromKB() uint64    { return uint64(b / KB) }
func (b ByteSize) FromMB() uint64    { return uint64(b / MB) }
func (b ByteSize) FromGB() uint64    { return uint64(b / GB) }
func (b ByteSize) FromTB() uint64    { return uint64(b / TB) }
func (b ByteSize) FromPB() uint64    { return uint64(b / PB) }
func (b ByteSize) FromEB() uint64    { return uint64(b / EB) }
func (b ByteSize) FromKiB() uint64   { return uint64(b / KiB) }
func (b ByteSize) FromMiB() uint64   { return uint64(b / MiB) }
func (b ByteSize) FromGiB() uint64   { return uint64(b / GiB) }
func (b ByteSize) FromTiB() uint64   { return uint64(b / TiB) }
func (b ByteSize) FromPiB() uint64   { return uint64(b / PiB) }
func (b ByteSize) FromEiB() uint64   { return uint64(b / EiB) }
