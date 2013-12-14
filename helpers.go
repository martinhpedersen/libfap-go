package fap

import "C"

func goUnsignedIntPtr(c_uint *_Ctype_uint) *uint {
	if c_uint == nil {
		return nil
	}
	v := uint(*c_uint)

	return &v
}

func goFloat64Ptr(c_double *_Ctype_double) *float64 {
	if c_double == nil {
		return nil
	}
	v := float64(*c_double)

	return &v
}

func goBoolPtr(c_short *_Ctype_short) *bool {
	if c_short == nil {
		return nil
	}
	v := uint(*c_short) == 1

	return &v
}

func goString(c_str *_Ctype_char) string {
	if c_str == nil {
		return ""
	}

	return C.GoString(c_str)
}
