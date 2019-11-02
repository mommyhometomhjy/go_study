package vm

import (
	"aliexpress/model"
)

type OrderViewModel struct {
	BaseViewModel
	Orders []model.Order
}
type OrderViewModelOp struct{}

func (OrderViewModelOp) OrderGetIndexVM() OrderViewModel {
	orders := model.GetOrders()
	v := OrderViewModel{
		BaseViewModel{Title: "订单列表"},
		orders,
	}
	// fmt.Println(orders[0])
	return v
}

func (OrderViewModelOp) OrderGetNewVM() OrderViewModel {

	v := OrderViewModel{
		BaseViewModel{Title: "新建订单"},
		nil,
	}
	// fmt.Println(orders[0])
	return v
}
