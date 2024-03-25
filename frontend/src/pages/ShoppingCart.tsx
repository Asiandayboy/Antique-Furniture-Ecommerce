import CartItem from "../components/CartItem";
import Navbar from "../components/Navbar";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";
import { Link } from "react-router-dom";

export default function ShoppingCart() {

  const { cart } = useShoppingCartContext()


  function getSubtotal(): number {
    let total = 0;
    for (const listing of Object.values(cart)) {
      total += +listing.cost
    }
    return total
  }

  function sendCheckoutRequest() {
    /* 
      reducing the cart items into an array of the listingIDs to 
      satisfy input of API endpoint
    */
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
      <div className="shopping-cart_wrapper">
        <h1>Shopping Cart</h1>
        <div className="cart_wrapper">
          <div className="cart-items_wrapper">
            {
              Object.entries(cart).map(([key, item]) => (
                <CartItem {...item} key={key}/>
              ))
            }
          </div>
          <div className="cart-info_wrapper">
            <div>
              <div className="cart-totals">
                {
                  Object.values(cart).map((item) => (
                    <div>
                      <div>{item.title}</div>
                      <div>${item.cost}</div>
                    </div>
                  ))
                }
              </div>
              <div className="cart-subtotal">Subtotal: ${getSubtotal()}</div>
            </div>
            {
              Object.entries(cart).length > 0 &&
              <button 
                className="checkout_btn"
                onClick={sendCheckoutRequest}
              >Checkout</button>
              ||
              <div className="no-checkout_btn">
                There are no items in your cart
              </div>
            }
          </div>
        </div>
      </div>
    </>
  )
}