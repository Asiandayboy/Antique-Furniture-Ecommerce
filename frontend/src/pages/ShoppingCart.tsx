import Navbar from "../components/Navbar";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";
import { Link } from "react-router-dom";

export default function ShoppingCart() {

  const { cart } = useShoppingCartContext()

  function sendCheckoutRequest() {
    // reducing the cart items into an array of the listingIDs 
    const cartItems = []
    let total = 0
    for (const [key, listing] of Object.entries(cart)) {
      cartItems.push(key)
      total += +listing.cost
    }

    const body = {
      shoppingCart: cartItems,
      paymentInfo: {
        paymentMethod: "credit",
        amount: total,
        currency: "usd"
      },
    }


    fetch("http://localhost:3000/checkout", {
      method: "POST",
      headers: {
        "Content-Type": "text/html",
      },
      body: JSON.stringify(body),
      credentials: "include"
    })
    .then(async (res: Response) => {
      if (res.ok) {
        const redirectURL = await res.text()
        console.log("Redirect URL:", redirectURL)
        window.location.href = redirectURL
      } else {
        const msg = await res.text()
        throw new Error(msg)
      }
    })
    .catch((err: Error) => {
      console.error(err)
    })

  }


  return (
    <>
      <Navbar />
      <div>
        <h1>Checkout</h1>
        <div>
          <h2>Shopping Cart</h2>
          {
            Object.entries(cart).map(([key, item]) => (
              <div key={item.listingID}>
                <div>Title: {item?.title}</div>
                <div>Description: {item?.description}</div>
                <div>Cost: {item?.cost}</div>
                <div>Material: {item?.material}</div>
                <div>Style: {item?.style}</div>
                <div>Type: {item?.type}</div>
                <div>Condition: {item?.condition}</div>
                <div>Bought: {String(item?.bought)}</div>
                <div>ListingID: {item?.listingID}</div>
                <div>SellerID: {item?.userID}</div>
                <br />
              </div>
            ))
          }

          {
            Object.entries(cart).length > 0 &&
            <button onClick={sendCheckoutRequest}>Checkout</button>
            ||
            <div>There are no items in your cart</div>
          }
        </div>
      </div>
    </>
  )
}