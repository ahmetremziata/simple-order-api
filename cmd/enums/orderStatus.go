package enum

type OrderStatus int

var (
	Created     OrderStatus = 1
	Approved    OrderStatus = 2
	Transferred OrderStatus = 3
	Shipped     OrderStatus = 4
	Delivered   OrderStatus = 5
)
