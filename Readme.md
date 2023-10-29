# Required payloads

## Signup
#### localhost:8081/api/v1/auth/signup POST
{
"first_name": "Tolu",
"last_name": "Thomas",
"email": "email@gmail.com",
"password": "secretpass",
"phone": "+123456789001"
}

## Login
#### localhost:8081/api/v1/auth/login POST
{
"email": "email@gmail.com",
"password": "secretpass",
}

## Create Order and Pay via paystack
#### localhost:8081/api/v1/create-order POST
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

## Get User Orders
#### localhost:8081/api/v1/orders GET