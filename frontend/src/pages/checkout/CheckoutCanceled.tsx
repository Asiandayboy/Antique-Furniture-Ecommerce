import { Link } from "react-router-dom";

export default function CheckoutCanceled() {
  return (
    <div className="checkout-canceled_page">
      <div>
        You have canceled your checkout.
        <Link to="/"> <span>Click here to return Home</span></Link>
      </div>
    </div>
  )
}