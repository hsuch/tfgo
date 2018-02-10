package TFGOServer

type HeldBy int
const (
	RED HeldBy= iota
	NEUTRAL
	BLUE
)

type ControlPoint struct {
	ID string

	Location Location
	Radius float64

	PayloadPath []Location // start, end
	PayloadLoc Location

	RedCount int
	BlueCount int
	CaptureStatus int
	HeldBy HeldBy
}

func (cp *ControlPoint) updateStatus() {

}