package hx

type Swap string

const (
	SwapInnerHTML   Swap = "innerHTML"
	SwapOuterHTML   Swap = "outerHTML"
	SwapTextContent Swap = "textContent"
	SwapBeforeBegin Swap = "beforebegin"
	SwapAfterBegin  Swap = "afterbegin"
	SwapBeforeEnd   Swap = "beforeend"
	SwapAfterEnd    Swap = "afterend"
	SwapDelete      Swap = "delete"
	SwapNone        Swap = "none"
)

func (s Swap) WithTransition() string {
	return string(s) + " transition:true"
}
