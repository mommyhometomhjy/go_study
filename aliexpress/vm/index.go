package vm

type IndexViewModel struct {
	BaseViewModel
	Words string
}
type IndexViewModelOp struct{}

func (IndexViewModelOp) GetVM() IndexViewModel {
	v := IndexViewModel{
		BaseViewModel{Title: "Homepage"},
		"你好",
	}
	return v
}
