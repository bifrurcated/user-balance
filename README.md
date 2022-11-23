# user-balance

#### Метод получения баланса пользователя. Принимает id пользователя.
```html
GET http://localhost:1234/api/v1/balance
Content-Type: application/json

{
    "user_id": 1
}
```
#### Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
```html
POST http://localhost:1234/api/v1/add-money
Content-Type: application/json

{
  "user_id": 1,
  "amount": 100
}
```
#### Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость.
```html
POST http://localhost:1234/api/v1/reserve
Content-Type: application/json
{
    "user_id": 1,
    "service_id": 3,
    "order_id": 2,
    "cost": 1500
}
```
#### Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму.
```html
POST http://localhost:1234/api/v1/reserve/profit
Content-Type: application/json

{
  "user_id": 1,
  "service_id": 3,
  "order_id": 2,
  "amount": 1500
}
```
#### Реализовать сценарий разрезервирования денег, если услугу применить не удалось.
```html
POST http://localhost:1234/api/v1/reserve/cancel
Content-Type: application/json

{
  "user_id": 1,
  "service_id": 1,
  "order_id": 1
}
```
#### Получения списка транзакций с пагинацией (укащывается значение последнего поля которое имеет тип по указаномму в сортировке значению, по умолчанию дата) и сортировкой по сумме (amount) и дате (datetime)
```html
GET http://localhost:1234/api/v1/history?sort_by=datetime&sort_order=desc&limit=3&last=199
Content-Type: application/json

{
  "user_id": 1
}
```
#### Перевод между пользователями
```html
POST http://localhost:1234/api/v1/transfer
Content-Type: application/json

{
  "sender_user_id": 2,
  "receiver_user_id": 1,
  "amount": 901
}
```