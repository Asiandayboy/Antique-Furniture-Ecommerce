import { useEffect } from "react"
import { Navigate } from "react-router-dom"


type Props = {
  setIsLoggedIn: React.Dispatch<React.SetStateAction<boolean>>
}



export default function Logout({ setIsLoggedIn }: Props) {
  const logout = async () => {
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
        throw new Error(msg)
      } else {
        setIsLoggedIn(false)

      }
    } catch (error) {
      console.log(error)
    }
  }

  logout();


  return (
    <Navigate to="/" />
  )
}