# 🛒 Core E-Commerce System (Go + Clean Architecture)

E-Commerce System Tech Stack : GO (FIber), PostgreSQL (database/sql)

## Features

### User 
- Register, Login, Logout, Refresh Token
- Refresh token stored securely in DB
- Middleware for authentication, authorization
- Role base Access

### Product Management
- Create, Update, Delete (Admin, Staff)
- Get all, Get single (Public)
- Assign multiple categories 

### Category Management
- Create, Update, Delete (Admin, Staff)
- Assign / remove from product

### Address Management
- Create, Update, Delete
- Get by ID or user
- Set Default
- Linked with order creation

### Cart System
- Add products to cart
- Get items in cart
- Remove items

### Order System
- Create order from cart items
- Copy user address snapshot
- Store order items (product, quantity, price)
- Auto calulate total price
- Deduct product stock 
- Transaction safe
