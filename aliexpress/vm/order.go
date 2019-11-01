package vm

import (
	"aliexpress/model"
)

type OrderViewModel struct {
	BaseViewModel
	Orders []model.Order
}
type OrderViewModelOp struct{}

func (OrderViewModelOp) GetVM() OrderViewModel {
	orders := model.GetOrders()
	v := OrderViewModel{
		BaseViewModel{Title: "订单列表"},
		orders,
	}
	// fmt.Println(orders[0])
	return v
}
