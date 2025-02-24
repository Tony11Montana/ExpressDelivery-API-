
create table Status_order(
id_status INT primary key auto_increment,
state varchar(20) not null 
);

create table Clients(
id_client int primary key auto_increment,
first_name varchar(100) not null,
last_name varchar(100) not null,
date_both date,
mobile_number varchar(100),
address varchar(100)
);

create table Employees(
id_employee int primary key auto_increment,
first_name varchar(100) not null,
last_name varchar(100) not null,
mobile_number varchar(100) not null,
salary int not null
);

CREATE TABLE Orders(
id_order INT primary key auto_increment,
id_client int not null,
id_employee int not null,
id_status int not null,
foreign key (id_client) references clients (id_client),  
foreign key (id_employee) references Employees (id_employee),
date_order DATE NOT NULL,
foreign key (id_status) references Status_order (id_status),
type_pay VARCHAR(20) default "bank card",
date_pay DATE not null,
description_order varchar(200),
price_delivery int not null
);

create table Warehouses(
id_warehouse INT primary key auto_increment,
id_employee INT not null,
foreign key (id_employee) references Employees (id_employee), 
name_warehouse varchar(50),
address_warehouse varchar(50)
);

create table Couriers(
id_courier int primary key auto_increment,
first_name varchar(100) not null,
last_name varchar(100) not null,
id_warehouse int not null,
foreign key (id_warehouse) references Warehouses(id_warehouse)
);

create table Products(
id_product int primary key auto_increment,
id_warehouse int not null,
foreign key (id_warehouse) references Warehouses (id_warehouse),
name_product VARCHAR(100) NOT NULL,
count_warehouse int check(count_warehouse > 0) not null,
description_product varchar(100)
);

create table Info_orders(
id_info_order int primary key auto_increment,
id_courier int not null,
id_order int not null,
id_product int not null,
foreign key (id_courier) references Couriers (id_courier),
foreign key (id_order) references Orders (id_order),
foreign key (id_product) references Products (id_product),
count_product INT CHECK(count_product > 0) not NULL,
price INT CHECK(price > 0) not NULL,
date_create date not null
);
