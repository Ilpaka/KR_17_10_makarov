1. СОздать структуру products с полями:
id;
name;
price;
in_stock;

2.Реализовать Rest API:
POST/products - добавление продукта
GET/produscts - спписок + фильтр min_price, max_price, in_stock
GET/products{:id}
DELTE/products{:id}