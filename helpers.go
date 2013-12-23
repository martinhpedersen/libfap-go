package fap

import (
	"C"
)

func goUnsignedInt(c_uint *_Ctype_uint) uint {
	if c_uint == nil {
		return 0
	}

	return uint(*c_uint)
}

func goFloat64(c_double *_Ctype_double) float64 {
	if c_double == nil {
		return 0
	}

	return float64(*c_double)
}

func goBool(c_short *_Ctype_short) bool {
	if c_short == nil {
		return false
	}

	return uint(*c_short) == 1
}

func goString(c_str *_Ctype_char) string {
	if c_str == nil {
		return ""
	}

	return C.GoString(c_str)
}
