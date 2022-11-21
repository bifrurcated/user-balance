# user-balance

#### Метод получения баланса пользователя. Принимает id пользователя.
```html
GET http://localhost:12345/api/v1/balance
Content-Type: application/json

{
    "user_id": 1
}
```
#### Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
```html
POST http://localhost:12345/api/v1/add-money
Content-Type: application/json

{
  "user_id": 1,
  "amount": 10
}
```