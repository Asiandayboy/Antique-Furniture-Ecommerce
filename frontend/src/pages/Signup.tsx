import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"

type SignupInfo = {
  username: string,
  email: string,
  password: string,
  confirm: string
}



function validateSignupInfo(data: SignupInfo): boolean {
  if (data.username == "") {
    return false
  }
  if (data.email == "") {
    return false
  }
  if (data.password == "") {
    return false
  }
  if (data.confirm == "") {
    return false
  }

  return true
}


export default function Signup() {
  const navigate = useNavigate()

  const [signupInfo, setSignupInfo] = useState<SignupInfo>({
    username: '',
    email: '',
    password: '',
    confirm: ''
  })

  const [success, setSuccess] = useState<boolean>(false)
  const [resMsg, setResMsg] = useState<string>("")

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    console.log("form submited:", signupInfo)

    const validated = validateSignupInfo(signupInfo)

    if (!validated) {
      setResMsg("Fields cannot be blank")
      return
    }

    try {
      const res = await fetch("http://localhost:3000/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(signupInfo)
      })

      if (!res.ok) {
        const msg = await res.text()
        setResMsg(msg)
        throw new Error(msg || "Failed to sign up!")
      } else {
        const msg = await res.text()
        console.log("signup:", msg)
        setSuccess(true)
        navigate("/signup-success")

      }
    } catch (error) {
      console.log(error)
    }

  }


  return (
    <main>
      <div className="signup_wrapper">
        <form className="signup-form" onSubmit={onSubmit}>
          <h1>Sign Up</h1>
          <div className="signup-user">
            <label htmlFor="username">Username</label>
            <input placeholder="Username" onChange={
              (e) => setSignupInfo({...signupInfo, username: e.target.value})  
            } type="text" name="username" id="username" />
          </div>
          <div className="signup-email">
            <label htmlFor="email">Email</label>
            <input placeholder="Email" onChange={
              (e) => setSignupInfo({...signupInfo, email: e.target.value})
            } type="email" name="email" id="email" />
          </div>
          <div className="signup-pass">
            <label htmlFor="password">Password</label>
            <input placeholder="Password" onChange={
              (e) => setSignupInfo({...signupInfo, password: e.target.value})
            }  type="password" name="password" id="password" />
          </div>
          <div className="signup-confirm">
            <label htmlFor="confirm">Confirm Password</label>
            <input placeholder="Confirm password" onChange={
              (e) => setSignupInfo({...signupInfo, confirm: e.target.value})
            }  type="password" name="confirm" id="confirm" />
          </div>
          
          <button type="submit" name="submit">Signup</button>
        </form>
        <div>
          {success && "Successfully signed up!" || resMsg}
        </div>
        <div>
          Have an account? <Link to="/login">Log in</Link>
        </div>
      </div>
    </main>
  )
}