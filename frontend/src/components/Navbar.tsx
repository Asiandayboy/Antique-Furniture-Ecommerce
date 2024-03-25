import { Link, useNavigate } from "react-router-dom"
import { useAuthContext } from "../contexts/authContext"
import { useShoppingCartContext } from "../contexts/shoppingCartContext"
import CartIcon from "../assets/CartIcon"
import HamburgerIcon from "../assets/HamburgerIcon"

type Props = {
  // isLoggedIn: boolean
}






export default function Navbar({  }: Props) {
  const navigate = useNavigate()

  const authState = useAuthContext()
  const shoppingCart = useShoppingCartContext()

  async function onLogoutClick(e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) {
    try {
      const res = await fetch("http://localhost:3000/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include"
      })

      if (!res.ok) {
        const msg = await res.text()

        if (res.status == 401) { // unauthorized
          throw new Error(msg)
        }
      } else {
        console.log("You have logged out!")
        authState.setAuth(false)
        localStorage.setItem("isLoggedIn", "false")
        navigate("/")
      }
    } catch (error) {
      console.error(error)

      authState.setAuth(false)
      localStorage.setItem("isLoggedIn", "false")
      navigate("/")

    }
  }


  return (
    <nav>
      <ol>
        <HamburgerIcon />
        <li className="nav-1">
          <Link to="/">Home</Link>
        </li>
        <div className="nav-2">
          <li><Link to="/list">List a furniture</Link></li>
          <li><Link to="/market">Market</Link></li>
          {!authState.auth &&
            <>
              <li>
                <Link to="/login">Login</Link>
              </li>
              <li className="signup-link">
                <button className="signup_btn"><Link to="/signup">Signup</Link></button>
              </li>
            </>
            ||
            <li>
              <div className="acc-dropdown">
                <Link className="acc-dropdown_btn" to="/dashboard">Account</Link>
                <div className="acc-dropdown-content">
                  <div className="acc-balance">Balance: $666.000</div>
                  <Link to="/dashboard">Account</Link>
                  <Link to="/dashboard/purchase-history">Purchase History</Link>
                  <Link to="/dashboard/furniture-listings">My Furniture Listings</Link>
                  <Link to="/dashboard/addresses">My Addresses</Link>
                  <Link to="" onClick={
                    (e) => onLogoutClick(e)
                  } className="logout-link">Logout</Link>
                </div>
              </div>
            </li>
          }
        </div>
        <li className="cart-link">
          <Link to="/shopping-cart">
            <CartIcon />
            <div>Cart ({Object.keys(shoppingCart.cart).length})</div>
          </Link>
        </li>
      </ol>
    </nav>
  )
}