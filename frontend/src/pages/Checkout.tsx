import { useEffect } from "react";
import Navbar from "../components/Navbar";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";

export default function Checkout() {

  const { cart } = useShoppingCartContext()

  function sendCheckoutRequest() {
    // reducing the cart items into an array of the listingIDs 
    const cartItems = []
    for (const key in cart) {
      cartItems.push(key)
    }

    // stub
    const body = {
      shoppingCart: cartItems,
      paymentInfo: {
        paymentMethod: "credit",
        amount: 666,
        currency: "usd"
      },
      shippingAddress: {
        state: "RI",
        city: "Providence",
        street: "First Test Checkout St.",
        zipCode: "02907"
      } 
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
        // const redirectURL = res.url
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


  useEffect(() => {

  })






  return (
    <>
      <Navbar />
      <div>
        <form>
          div

          <button type="submit">Checkout</button>
        </form>
      </div>
    </>
  )
}