# Clothing Store API — руководство пользователя

## 1. Подготовка к работе

Открой PowerShell и перейди в папку проекта:

# cd C:\Users\mozad\Desktop\clothing-store

Укажи адрес сервера:

# $base = "http://localhost:8080"

Все дальнейшие запросы используют переменную $base.

Для отправки запросов используется Invoke-RestMethod.

    Обозначения IDПОЛЬЗОВАТЕЛЯ, IDТОВАРА, IDЗАКАЗА и т. д. — места, куда нужно подставить настоящие UUID из ответов сервера.

## 2. Role — роли пользователей

Роли определяют тип пользователя. В базе уже имеются четыре начальные роли:

Название \  ID
    guest   11111111-1111-1111-1111-111111111111
    client  22222222-2222-2222-2222-222222222222
    manager 33333333-3333-3333-3333-333333333333
    admin   44444444-4444-4444-4444-444444444444

# GET — получить все роли

# Invoke-RestMethod -Method GET -Uri "$base/role"

Результат: список всех ролей с их UUID и названиями.

# POST — создать роль

# Invoke-RestMethod -Method POST -Uri "$base/role" -ContentType "application/json" -Body '{"name":"НАЗВАНИЕ_РОЛИ"}'

Пример:

    Invoke-RestMethod -Method POST -Uri "$base/role" -ContentType "application/json" -Body '{"name":"test_role"}'

Результат: сервер создаст роль и вернёт её UUID и название.

# PUT — изменить роль

# Invoke-RestMethod -Method PUT -Uri "$base/role/IDРОЛИ" -ContentType "application/json" -Body '{"name":"НОВОЕ_НАЗВАНИЕ"}'

Пример:

    Invoke-RestMethod -Method PUT -Uri "$base/role/IDРОЛИ" -ContentType "application/json" -Body '{"name":"updated_role"}'

Результат: название роли изменится.

# DELETE — удалить роль

# Invoke-RestMethod -Method DELETE -Uri "$base/role/IDРОЛИ"

Результат: роль будет удалена.

## 3. User — пользователи

Каждый пользователь должен иметь существующую роль. При создании передаётся role_id.

# GET — получить всех пользователей

# Invoke-RestMethod -Method GET -Uri "$base/user?limit=100&offset=0"

Результат: список пользователей с именами, электронной почтой, телефонами, UUID ролей и датами регистрации.

# POST — создать пользователя

# Invoke-RestMethod -Method POST -Uri "$base/user" -ContentType "application/json" -Body '{"name":"ИМЯ","email":"EMAIL","phone":"ТЕЛЕФОН","role_id":"IDРОЛИ"}'

Пример:

    Invoke-RestMethod -Method POST -Uri "$base/user" -ContentType "application/json" -Body '{"name":"John Smith","email":"john.smith@example.com","phone":"+12025550101","role_id":"22222222-2222-2222-2222-222222222222"}'

Результат: сервер создаст пользователя, назначит ему UUID и дату регистрации.

# GET — получить пользователей определённой роли

# Invoke-RestMethod -Method GET -Uri "$base/user/role/IDРОЛИ"

Пример:

    Invoke-RestMethod -Method GET -Uri "$base/user/role/22222222-2222-2222-2222-222222222222"

Результат: список пользователей, привязанных к указанной роли.

# PUT — изменить пользователя

# Invoke-RestMethod -Method PUT -Uri "$base/user/IDПОЛЬЗОВАТЕЛЯ" -ContentType "application/json" -Body '{"name":"НОВОЕ_ИМЯ","email":"EMAIL","phone":"ТЕЛЕФОН","role_id":"IDРОЛИ"}'

Пример:

    Invoke-RestMethod -Method PUT -Uri "$base/user/IDПОЛЬЗОВАТЕЛЯ" -ContentType "application/json" -Body '{"name":"John Smith Updated","email":"john.smith@example.com","phone":"+12025550999","role_id":"22222222-2222-2222-2222-222222222222"}'

Результат: данные пользователя обновятся.

# DELETE — удалить пользователя

# Invoke-RestMethod -Method DELETE -Uri "$base/user/IDПОЛЬЗОВАТЕЛЯ"

Результат: пользователь будет удалён.

## 4. Product — товары

Товар содержит название, тип, пол, размер, цвет, описание, бренд, сезон и цену.

# GET — получить все товары

# Invoke-RestMethod -Method GET -Uri "$base/product"

Результат: список всех товаров.

# POST — создать товар

# Invoke-RestMethod -Method POST -Uri "$base/product" -ContentType "application/json" -Body '{"name":"НАЗВАНИЕ","type":"ТИП","gender":"ПОЛ","size":"РАЗМЕР","color":"ЦВЕТ","description":"ОПИСАНИЕ","brand":"БРЕНД","season":"СЕЗОН","price":ЦЕНА}'

Пример:

    Invoke-RestMethod -Method POST -Uri "$base/product" -ContentType "application/json" -Body '{"name":"Classic Black T-Shirt","type":"T-Shirt","gender":"Male","size":"M","color":"Black", "description":"Classic cotton t-shirt for everyday wear","brand":"Nike","season":"Summer","price":29.99}'

Результат: сервер создаст товар и вернёт его UUID и характеристики.

# GET — получить один товар

# Invoke-RestMethod -Method GET -Uri "$base/product/IDТОВАРА"

Результат: информация об одном товаре.

# GET — фильтрация товаров

По типу:

    Invoke-RestMethod -Method GET -Uri "$base/product?type=T-Shirt"

По полу:

    Invoke-RestMethod -Method GET -Uri "$base/product?gender=Male"

По размеру:

    Invoke-RestMethod -Method GET -Uri "$base/product?size=M"

По цвету:

    Invoke-RestMethod -Method GET -Uri "$base/product?color=Black"

По бренду:

    Invoke-RestMethod -Method GET -Uri "$base/product?brand=Nike"

По сезону:

    Invoke-RestMethod -Method GET -Uri "$base/product?season=Summer"

Результат: список товаров, соответствующих указанному фильтру.

# GET — пагинация

# Invoke-RestMethod -Method GET -Uri "$base/product?limit=КОЛИЧЕСТВО&offset=СМЕЩЕНИЕ"

Пример:

    Invoke-RestMethod -Method GET -Uri "$base/product?limit=2&offset=0"

Результат: не более двух товаров, начиная с первого.

# PUT — изменить товар

# Invoke-RestMethod -Method PUT -Uri "$base/product/IDТОВАРА" -ContentType "application/json" -Body '{"name":"НОВОЕ_НАЗВАНИЕ","type":"ТИП","gender":"ПОЛ","size":"РАЗМЕР","color":"ЦВЕТ","description":"ОПИСАНИЕ","brand":"БРЕНД","season":"СЕЗОН","price":ЦЕНА}'

Пример:

    Invoke-RestMethod -Method PUT -Uri "$base/product/IDТОВАРА" -ContentType "application/json" -Body '{"name":"Classic Black T-Shirt Updated","type":"T-Shirt","gender":"Male","size":"XL","color":"White","description":"Updated cotton t-shirt","brand":"Nike","season":"Summer","price":34.99}'

Результат: характеристики товара обновятся.

# DELETE — удалить товар

# Invoke-RestMethod -Method DELETE -Uri "$base/product/IDТОВАРА"

Результат: товар будет удалён.

## 5. Cart — корзина

Корзина связывает пользователя с товаром и хранит количество единиц товара.

# POST — добавить товар в корзину

# Invoke-RestMethod -Method POST -Uri "$base/cart" -ContentType "application/json" -Body '{"user_id":"IDПОЛЬЗОВАТЕЛЯ","product_id":"IDТОВАРА","count":КОЛИЧЕСТВО}'

Пример:

    Invoke-RestMethod -Method POST -Uri "$base/cart" -ContentType "application/json" -Body '{"user_id":"IDПОЛЬЗОВАТЕЛЯ","product_id":"IDТОВАРА","count":2}'

Результат: товар будет добавлен в корзину. Сервер вернёт ID записи корзины.

# GET — получить корзину пользователя

# Invoke-RestMethod -Method GET -Uri "$base/cart/IDПОЛЬЗОВАТЕЛЯ"

Результат: список товаров в корзине указанного пользователя.

# PUT — изменить количество товара

# Invoke-RestMethod -Method PUT -Uri "$base/cart/IDЗАПИСИ_КОРЗИНЫ" -ContentType "application/json" -Body '{"count":НОВОЕ_КОЛИЧЕСТВО}'

Пример:

    Invoke-RestMethod -Method PUT -Uri "$base/cart/IDЗАПИСИ_КОРЗИНЫ" -ContentType "application/json" -Body '{"count":5}'

Результат: количество товара в выбранной записи корзины изменится.

# DELETE — удалить запись корзины

# Invoke-RestMethod -Method DELETE -Uri "$base/cart/IDЗАПИСИ_КОРЗИНЫ"

## 6. Order — заказы

Заказ связывает покупателя с датой создания и текущим статусом.

Допустимые статусы:

    Принят
    В работе
    Собран
    Готов к выдаче
    Отменён

# POST — создать заказ

Пример:
    $order = Invoke-RestMethod -Method POST -Uri "$base/order" -ContentType "application/json" -Body '{"user_id":"b61f7da6-a996-4aca-b82a-80ab2cfcf68d"}'


Результат: сервер создаст заказ, присвоит ему UUID, дату и начальный статус Принят.

# GET — получить все заказы

Пример:
    Invoke-RestMethod -Method GET -Uri "$base/order?limit=100&offset=0"


Результат: список заказов.

# GET — получить один заказ

Пример:
    Invoke-RestMethod -Method GET -Uri "$base/order/$($order.id)"


Результат: информация о выбранном заказе.

# GET — фильтрация заказов по статусу

Пример:
    Invoke-RestMethod -Method GET -Uri "$base/order?status=%D0%9F%D1%80%D0%B8%D0%BD%D1%8F%D1%82"


Результат: список заказов со статусом Принят.

# GET — фильтрация заказов по дате

Пример:
    Invoke-RestMethod -Method GET -Uri "$base/order?from=2026-11-01T00:00:00&to=2026-11-30T23:59:59"

Результат: список заказов за указанный период.

# PUT — изменить статус заказа

Пример:
    Invoke-RestMethod -Method PUT -Uri "$base/order/$($order.id)" -ContentType "application/json" -Body '{"status":"В работе"}'

Результат: статус заказа изменится.

## 7. OrderItem — позиции заказа

Позиция заказа связывает конкретный товар с конкретным заказом и хранит его цену.

# POST — добавить товар в заказ

Пример:
    Invoke-RestMethod -Method POST -Uri "$base/order_item" -ContentType "application/json" -Body (@{order_id=$order.id; product_id="93c0d081-5ae1-4ae7-905a-387bbe1b11aa"; price=119.99} | ConvertTo-Json)

Результат: сервер создаст позицию заказа и вернёт её UUID.

# GET — получить позиции заказа

Пример:
    Invoke-RestMethod -Method GET -Uri "$base/order_item/$($order.id)"

Результат: список товаров, добавленных в указанный заказ.
