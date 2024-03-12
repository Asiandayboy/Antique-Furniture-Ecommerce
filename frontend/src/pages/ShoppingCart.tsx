import Navbar from "../components/Navbar";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";
import { Link } from "react-router-dom";

export default function ShoppingCart() {

  const { cart } = useShoppingCartContext()


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
          <Link to="/checkout"><button>Checkout</button></Link>
        </div>
      </div>
    </>
  )
}