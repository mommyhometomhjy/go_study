package vm

import (
	"aliexpress/model"
)

type OrderViewModel struct {
	BaseViewModel
	Orders []model.Order
	Order  model.Order
}
type OrderViewModelOp struct{}

func (OrderViewModelOp) OrderGetIndexVM() OrderViewModel {
	orders := model.GetOrders()
	v := OrderViewModel{
		BaseViewModel{Title: "订单列表"},
		orders,
		model.Order{},
	}
	// fmt.Println(orders[0])
	return v
}

func (OrderViewModelOp) OrderGetNewVM() OrderViewModel {

	v := OrderViewModel{
		BaseViewModel{Title: "新建订单"},
		nil,
		model.Order{},
	}
	// fmt.Println(orders[0])
	return v
}
func (OrderViewModelOp) OrderGetEditVM(id int) OrderViewModel {
	order := model.GetOrderById(id)
	v := OrderViewModel{
		BaseViewModel{Title: "编辑订单"},
		nil,
		order,
	}
	return v
}
