# Required payloads

## Signup
#### localhost:8081/api/v1/auth/signup POST
#### https://escrolla.onrender.com/api/v1/auth/signup POST
{
"first_name": "Tolu",
"last_name": "Thomas",
"email": "email@gmail.com",
"password": "secretpass",
"phone": "+123456789001"
}

## Login
#### localhost:8081/api/v1/auth/login POST
#### https://escrolla.onrender.com/api/v1/auth/login POST
{
"email": "email@gmail.com",
"password": "secretpass",
}

## Create Order and Pay via paystack
#### localhost:8081/api/v1/create-order POST
#### https://escrolla.onrender.com/api/v1/create-order POST
{
"buyer_phone": "123-456-7890",
"seller_phone": "987-654-3210",
"buyer_email": "buyer@example.com",
"seller_email": "seller@example.com",
"amount": 5000,
"description": "Product description with terns and condition",
"delivery_days": 5,
"user_type": "buyer"
}

## Verify Paystack Trnsaction
#### localhost:8081/api/v1/verify?reference=a98bedf1-83c3-41f5-ba2b-7e9cf9924e4f  GET
#### https://escrolla.onrender.com/api/v1/verify?reference=a98bedf1-83c3-41f5-ba2b-7e9cf9924e4f  GET

## Get User Orders
#### localhost:8081/api/v1/orders GET
#### Response
"data": {
"orders": [
{
"id": "a98bedf1-83c3-41f5-ba2b-7e9cf9924e4f",
"created_at": 1698551350,
"updated_at": 1698551894,
"deleted_at": null,
"user_id": "1",
"buyer_phone": "123-456-7890",
"seller_phone": "987-654-3210",
"buyer_email": "buyer@example.com",
"seller_email": "seller@example.com",
"amount": 5000,
"description": "Product description with terns and condition",
"delivery_days": 5,
"user_type": "buyer",
"order_status": "pending",
"payment_status": "paid",
"EscrowFee": 100
},
{
"id": "b5b8c8c3-9cee-45e2-8728-4d87f777b412",
"created_at": 1698549936,
"updated_at": 1698549936,
"deleted_at": null,
"user_id": "1",
"buyer_phone": "123-456-7890",
"seller_phone": "987-654-3210",
"buyer_email": "buyer@example.com",
"seller_email": "seller@example.com",
"amount": 5000,
"description": "Product description with terns and condition",
"delivery_days": 5,
"user_type": "buyer",
"order_status": "pending",
"payment_status": "",
"EscrowFee": 100
},
{
"id": "fd729cd8-6c37-415c-ad11-9085db8fa20c",
"created_at": 1698550806,
"updated_at": 1698550806,
"deleted_at": null,
"user_id": "1",
"buyer_phone": "123-456-7890",
"seller_phone": "987-654-3210",
"buyer_email": "buyer@example.com",
"seller_email": "seller@example.com",
"amount": 5000,
"description": "Product description with terns and condition",
"delivery_days": 5,
"user_type": "buyer",
"order_status": "pending",
"payment_status": "pending",
"EscrowFee": 100
}
]
},
"errors": "",
"message": "retrieved orders successfully",
"status": "OK"
}