package TFGOServer

type ControlPoint struct {
	ID string

	Location Location
	Radius float64

	PayloadPath []Location // start, end
	PayloadLoc Location

	RedCount int
	BlueCount int
	CaptureStatus int
}
