package fap

import (
	"C"
)

func goUnsignedInt(c_uint *C.uint) uint {
	if c_uint == nil {
		return 0
	}

	return uint(*c_uint)
}

func goFloat64(c_double *C.double) float64 {
	if c_double == nil {
		return 0
	}

	return float64(*c_double)
}

func goBool(c_short *C.short) bool {
	if c_short == nil {
		return false
	}

	return uint(*c_short) == 1
}

func goString(c_str *C.char) string {
	if c_str == nil {
		return ""
	}

	return C.GoString(c_str)
}
