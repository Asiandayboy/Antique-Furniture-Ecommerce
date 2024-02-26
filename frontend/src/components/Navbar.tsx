import { Link, useNavigate } from "react-router-dom"
import { useAuthContext } from "../contexts/authContext"


type Props = {
  // isLoggedIn: boolean
}






export default function Navbar({  }: Props) {
  const navigate = useNavigate()

  const authState = useAuthContext()

  async function onLogoutClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
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
        <li className="nav-1">
          <Link to="/">Home</Link>
        </li>
        <div className="nav-2">
          <li>
            <Link to="/market">Market</Link>
          </li>
          {!authState.auth &&
            <>
              <li>
                <Link to="/login">Login</Link>
              </li>
              <li className="signup-link">
                <button><Link to="/signup">Signup</Link></button>
              </li>
            </>
            ||
            <>
              <li>
                <button onClick={onLogoutClick}>Logout</button>
              </li>
              <li>
                <Link to="/dashboard">Account</Link>
              </li>
            </>
          }
          <li>Cart (0)</li>
        </div>
      </ol>
    </nav>
  )
}