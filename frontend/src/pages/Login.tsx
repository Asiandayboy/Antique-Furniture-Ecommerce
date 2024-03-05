import { Link, useNavigate } from "react-router-dom"
import { useState } from "react"

type LoginInfo = {
  username: string,
  password: string
}

type Props = {
  setIsLoggedIn: React.Dispatch<React.SetStateAction<boolean>>
}


export default function Login({ setIsLoggedIn }: Props) {
  const navigate = useNavigate()
  const [loginInfo, setLoginInfo] = useState<LoginInfo>({
    username: "",
    password: "",
  })

  const [resMsg, setResMsg] = useState<string>("")

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    try {
      const res = await fetch("http://localhost:3000/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify(loginInfo)
      })

      if (!res.ok) {
        const msg = await res.text();
        setResMsg(msg)
        throw new Error(msg || "Failed to log in")
      } else {
        setIsLoggedIn(true)

        const msg = await res.text();
        console.log(msg)
        navigate("/dashboard")
      }

    } catch (error) {
      console.log(error)
    }
  }


  return (
    <main>
      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="username">Username</label>
          <input onChange={
            (e) => setLoginInfo({...loginInfo, username: e.target.value})
          } type="text" name="username" id="username" />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input onChange={
            (e) => setLoginInfo({...loginInfo, password: e.target.value})
          } type="password" name="password" id="password" />
        </div>
        <button type="submit" name="submit">Login</button>
      </form>
      <div>{resMsg}</div>
      <div>
        Don't have an account? <Link to="/signup">Sign up</Link>
      </div>
    </main>
  )
}