package cache

import "strings" 

type ScaleMode int

const (
	SCALEMODE_RESIZE  = ScaleMode(iota)
	SCALEMODE_FIT     = ScaleMode(iota)
	//SCALEMODE_FILL  = ScaleMode(iota)
	//SCALEMODE_CROP  = ScaleMode(iota)
	SCALEMODE_DEFAULT = SCALEMODE_RESIZE
)

func (m *ScaleMode) FromString(s string) {
	s = strings.ToLower(s)
	switch s {
	case "resize":
		*m = SCALEMODE_RESIZE
		return
	case "fit":
		*m = SCALEMODE_FIT
		return
	}
	*m = SCALEMODE_DEFAULT
	return
}
