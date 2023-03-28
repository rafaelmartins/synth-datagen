package ctypes

import "fmt"

func printformat(i interface{}, hex bool) string {
	return fmt.Sprint(i)
}

func stringformat(i interface{}, hex bool) string {
	return fmt.Sprintf("%q", i)
}

func valueformat(i interface{}, hex bool) string {
	if !hex {
		return printformat(i, false)
	}

	switch v := i.(type) {
	case int8:
		return fmt.Sprintf("0x%02x", uint8(v))
	case uint8:
		return fmt.Sprintf("0x%02x", v)
	case int16:
		return fmt.Sprintf("0x%04x", uint16(v))
	case uint16:
		return fmt.Sprintf("0x%04x", v)
	case int32:
		return fmt.Sprintf("0x%08x", uint32(v))
	case uint32:
		return fmt.Sprintf("0x%08x", v)
	case int64:
		return fmt.Sprintf("0x%016x", uint64(v))
	case uint64:
		return fmt.Sprintf("0x%016x", v)
	case int:
		return fmt.Sprintf("0x%08x", uint32(v))
	case uint:
		return fmt.Sprintf("0x%08x", uint32(v))
	}

	return "0"
}
