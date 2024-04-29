# How to run the app

*We need to set up some constants first...*

### Installing Stripe
- Make a Stripe Account and install the Stripe CLI
- Store your Stripe secret key as an environment system variable named `STRIPE_TEST_KEY`

### Installing MongoDB
- Install MongoDB (I have MongoDB Compass installed as well, which is the GUI)
- Connect to MongoDB with this connection string: `mongodb://localhost:27017`
- Create a new database called `AntiqueFurnitureProject`
- Create four collections in the database you just created named:
    - `listings`
    - `receipts`
    - `shippingAddresses`
    - `users`
- Import each of the database_dump JSON files in its respective collection

___

Once that is done, you can clone the repository into your local environment, and open up two terminals: one for the frontend and backend. 
- Move into the frontend directory in terminal 1, and run `npm run dev`
- Move into the backend directory in terminal 2, and run `go run main.go`

Now that the app has started, just go to the URL of the React app. It should just print it out in terminal 1.


#### Valid Logins
- username: `varus`; password: `varuslover`
- username: `Asiandayboy`; password: `password`



---
##### Note

*The app doesn't have a lot of data for furniture listings, but more can be made by creating a new furniture listing.*