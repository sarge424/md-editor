package kb

import (
	"fmt"
	"unicode"
)

type State struct {
	Ctrl  int
	Alt   int
	Shift int
}

type Shortcut struct {
	Ctrl  bool
	Alt   bool
	Shift bool
	Caps  bool

	Code int
	Char rune
}

type Keystroke rune

type Keyboarder interface {
	HandleKeystroke(k Keystroke)
	HandleShortcut(s Shortcut)
}

func (kb *State) Emit(scancode int, rn rune) Shortcut {
	s := kb.generate()
	s.Code = scancode
	s.Char = rn

	return s
}

func (kb State) generate() Shortcut {
	var s Shortcut
	s.Ctrl = kb.Ctrl > 0
	s.Alt = kb.Alt > 0
	s.Shift = kb.Shift > 0

	return s
}

func (ks Keystroke) Std() rune {
	r := rune(ks)
	if unicode.IsLetter(r) {
		return unicode.ToUpper(r)
	}

	return r
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

	if s.Char == 10 {
		ans += "ENTER"
	} else if s.Char == rune('\t') {
		ans += "TAB"
	} else if s.Code == 14 {
		ans += "BCKSP"
	} else if s.Char != 0 {
		ans += string(s.Char)
	} else {
		ans += fmt.Sprint(s.Code)
	}

	return ans
}

func (k Keystroke) String() string {
	return string(k)
}
