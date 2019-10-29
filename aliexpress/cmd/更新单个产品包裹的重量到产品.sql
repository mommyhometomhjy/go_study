update goods
set goods_weight = (select we from (select avg(weight) as we,goods_no 
from (select 
	order_shipping_weight as weight,
	goods_no
from orders 
join order_details on orders.id = order_details.order_id
join goods on order_details.goods_id = goods.id
where order_shipping_weight >0
group by order_shipping_no 
having count(order_no)=1)
group by goods_no) as t2 where t2.goods_no = goods.goods_no limit 1)
where exists(select * from (select avg(weight) as we,goods_no 
from (select 
	order_shipping_weight as weight,
	goods_no
from orders 
join order_details on orders.id = order_details.order_id
join goods on order_details.goods_id = goods.id
where order_shipping_weight >0
group by order_shipping_no 
having count(order_no)=1)
group by goods_no) as t2 where t2.goods_no = goods.goods_no )