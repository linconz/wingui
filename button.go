package wingui

import "log"

type Button struct {
	WindowBase
	OnClicked func()
}

func (b *Button) WndProc(msg uint32, wParam, lParam uintptr) uintptr {
	log.Println("btn wnd proc", b.hwnd, msg, wParam, lParam)
	if b.OnClicked != nil {
		b.OnClicked()
	}
	b.AsWindowBase().WndProc(msg, wParam, lParam)
	return 0
}
