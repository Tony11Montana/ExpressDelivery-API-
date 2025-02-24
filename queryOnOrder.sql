SELECT Orders_new.id_client, concat(Clients.first_name," " ,Clients.last_name) as Client,
Orders_new.price_delivery, Orders_new.type_pay, Orders_new.date_pay, Orders_new.date_order,
Orders_new.CourierName, Orders_new.count_product, Orders_new.price,
Orders_new.name_product, Orders_new.count_warehouse
from (
select Orders.id_client, Orders.price_delivery, Orders.type_pay, Orders.date_pay, Orders.date_order,
info_cour.CourierName, info_cour.count_product, info_cour.price,
info_cour.name_product, info_cour.count_warehouse
from Orders inner join ( 
select concat(Couriers.first_name, " " , Couriers.last_name) as CourierName,
Info_orders_new.count_product, Info_orders_new.price, Info_orders_new.id_order,
Info_orders_new.name_product, Info_orders_new.count_warehouse
from (select Info_orders.id_courier, Info_orders.count_product, Info_orders.price, Info_orders.id_order,
Products.name_product, Products.count_warehouse
from Info_orders INNER join Products
) as Info_orders_new INNER JOIN Couriers ON Info_orders_new.id_courier = Couriers.id_courier
) as info_cour ON Orders.id_order = info_cour.id_order
) as Orders_new INNER JOIN Clients ON
Orders_new.id_client = Clients.id_client;
