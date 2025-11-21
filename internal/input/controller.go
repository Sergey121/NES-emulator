package input

const (
	ButtonA      = 1 << 0
	ButtonB      = 1 << 1
	ButtonSelect = 1 << 2
	ButtonStart  = 1 << 3
	ButtonUp     = 1 << 4
	ButtonDown   = 1 << 5
	ButtonLeft   = 1 << 6
	ButtonRight  = 1 << 7
)

type Controller struct {
	buttons byte // Current state of buttons (set by emulator host)
	strobe  byte // Strobe mode (0 or 1)
	state   byte // Internal shift register for serial reading
}

func NewController() *Controller {
	return &Controller{}
}

// SetButtons updates the current state of the buttons from the input source (keyboard/gamepad)
func (c *Controller) SetButtons(buttons byte) {
	c.buttons = buttons
}

// Write handles writing to the controller port (usually $4016)
// Writing 1 sets strobe mode (continuously reloading state).
// Writing 0 clears strobe mode.
func (c *Controller) Write(data byte) {
	c.strobe = data & 1
	if c.strobe == 1 {
		c.state = c.buttons
	}
}

// Read returns the next button state bit.
// Standard NES controller returns 1 bit at a time: A, B, Select, Start, Up, Down, Left, Right.
// After 8 reads, it usually returns 1s.
func (c *Controller) Read() byte {
	if c.strobe == 1 {
		// If strobe is high, we always read the current state of the A button (bit 0)
		return c.buttons & 1
	}

	// Get the low bit of the current state
	value := c.state & 1
	// Shift the state to prepare the next button
	c.state >>= 1
	// Standard NES controllers return 1s after the 8 buttons are read.
	// We simulate this by shifting in 1s.
	c.state |= 0x80

	return value
}
