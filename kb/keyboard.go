package kb

type State struct {
	Ctrl  int
	Alt   int
	Shift int
	Caps  bool
}

type Shortcut struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	Caps  bool

	Code int
	Key  rune
}

func (kb *State) Emit(scancode int, rn rune) (s Shortcut, valid bool) {
	s = kb.generate()
	s.Code = scancode
	s.Key = rn

	valid = true
	if s.Key == 0 && !(scancode == 10 || scancode == 14) {
		valid = false
	}

	return s, valid
}

func (kb State) generate() Shortcut {
	var s Shortcut
	s.Ctrl = kb.Ctrl > 0
	s.Alt = kb.Alt > 0
	s.Shift = kb.Shift > 0
	s.Caps = kb.Caps

	return s
}

func (kb *State) HandleUnderflow() {
	kb.Ctrl = max(kb.Ctrl, 0)
	kb.Alt = max(kb.Alt, 0)
	kb.Shift = max(kb.Shift, 0)
}

func (s Shortcut) String() string {
	ans := ""
	if s.Ctrl {
		ans += "CTL "
	}
	if s.Alt {
		ans += "ALT "
	}
	if s.Shift {
		ans += "SHF "
	}
	if s.Caps {
		ans += "CAP "
	}

	if s.Key == 10 {
		ans += "ENTER"
	} else if s.Key == rune('\t') {
		ans += "TAB"
	} else if s.Code == 14 {
		ans += "BCKSP"
	} else {
		ans += string(s.Key)
	}

	return ans
}
