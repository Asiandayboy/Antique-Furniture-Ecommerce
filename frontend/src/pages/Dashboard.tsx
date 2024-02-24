import { useEffect, useState } from "react"
import Navbar from "../components/Navbar"

type Props = {
  // isLoggedIn: boolean
}


type Info = {
  UserID: string,
  username: string,
  email: string,
  password: string,
  phone: string,
  sessionId: string,
  balance: string
}


export default function Dashboard({  }: Props) {
  const [userInfo, setUserInfo] = useState<Info | null>(null)

  async function onClick(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    try {
      const res = await fetch("http://localhost:3000/account", {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        },
        mode: "cors",
        credentials: "include",
      })

      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      } else {
        const data = await res.json()
        setUserInfo({...data})
      }

    } catch (error) {
      console.log(error)
    }
  }

  useEffect(() => {
    
  })



  return (
    <>
      <Navbar />
      <main>
        <div>Dashboard</div>
        <button onClick={(e) => onClick(e)}>Test button to get account info</button>
        <div>
          <div>UserID: {userInfo?.UserID}</div>
          <div>Username: {userInfo?.username}</div>
          <div>Email: {userInfo?.email}</div>
          <div>Password: {userInfo?.password}</div>
          <div>Phone: {userInfo?.phone}</div>
          <div>SessionID: {userInfo?.sessionId}</div>
          <div>Balance: {userInfo?.balance}</div>
        </div>
      </main>
    </>
  )
}