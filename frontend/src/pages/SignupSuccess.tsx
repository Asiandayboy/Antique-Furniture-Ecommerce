import { Link } from "react-router-dom"

export default function SignupSuccess() {
  return (
    <div className="signup-success_page">
      <div>
        You have successfully signed up! You can now 
        <Link to="/login"> <span>login</span></Link>
      </div>
    </div>
  )
}