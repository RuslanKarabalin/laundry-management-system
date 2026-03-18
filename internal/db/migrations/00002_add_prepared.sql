-- +goose Up
insert into appliances(name, type) values('Первая стиралка', 'washing_machine');
insert into appliances(name, type) values('Вторая стиралка', 'washing_machine');
insert into appliances(name, type) values('Третья стиралка', 'washing_machine');
insert into appliances(name, type) values('Четвертая стиралка', 'washing_machine');
insert into appliances(name, type) values('Пятая стиралка', 'washing_machine');

insert into appliances(name, type) values('Первая сушилка', 'tumble_dryer');
insert into appliances(name, type) values('Вторая сушилка', 'tumble_dryer');
insert into appliances(name, type) values('Третья сушилка', 'tumble_dryer');
insert into appliances(name, type) values('Четвертая сушилка', 'tumble_dryer');
insert into appliances(name, type) values('Пятая сушилка', 'tumble_dryer');

-- +goose Down
truncate table appliances;
