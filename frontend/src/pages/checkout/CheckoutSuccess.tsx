import { Link } from "react-router-dom";

export default function CheckoutSuccess() {
  return (
    <div className="checkout-success_page">
      <div>
        Your order has successfully been purchased. 
        You can now close this window or
        <Link to="/"> <span>click here to return Home</span></Link>
      </div>
    </div>
  )
}