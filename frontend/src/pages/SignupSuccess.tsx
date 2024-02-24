import { Link } from "react-router-dom"

export default function SignupSuccess() {
  return (
    <div>
      You have successfully signed up! You can now <Link to="/login">login</Link>
    </div>
  )
}