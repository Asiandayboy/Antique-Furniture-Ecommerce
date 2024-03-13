import { Link } from "react-router-dom";

export default function CheckoutSuccess() {
  return (
    <div>
      Your order has successfully been purchased. 
      You can now close this window.
      <Link to="/"><button>Click here to return Home</button></Link>
    </div>
  )
}