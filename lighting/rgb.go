package lighting

type RGB struct {
	R, G, B uint8
}

func (c RGB) ToUint32() (color uint32) {
	color |= 0xff000000
	color |= uint32(c.R) << 16
	color |= uint32(c.G) << 8
	color |= uint32(c.B)

	return
}
