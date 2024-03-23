import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import Navbar from "../components/Navbar"
import HouseIcon from "../assets/HouseIcon"
import { useAccountDataContext } from "../contexts/accountDataContext"

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
  const { userData, setUserData } = useAccountDataContext()

  const navigate = useNavigate()

  async function getAccountData() {
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
        console.log("account data:", data)
        setUserData({...data})
      }

    } catch (error) {
      console.log(error)
    }
  }

  useEffect(() => {
    getAccountData()
  }, [])



  return (
    <>
      <Navbar />
      <main>
        <div className="dashboard_wrapper">
          <h1>Dashboard</h1>
          <div className="dashboard-grid">
            <div className="dashboard-grid-item"
              onClick={() => navigate("/dashboard/purchase-history")}
            >
              <HouseIcon />
              <div>
                <div className="header">Purchase History</div>
                <div className="subtext">
                  View all of your purchases
                </div>
              </div>
            </div>
            <div className="dashboard-grid-item"
              onClick={() => navigate("/dashboard/login-security")}
            >
              <HouseIcon />
              <div>
                <div className="header">Login & Security</div>
                <div className="subtext">
                  Edit your password, email, and phone number
                </div>
              </div>
            </div>
            <div className="dashboard-grid-item"
              onClick={() => navigate("/dashboard/addresses")}
            >
              <HouseIcon />
              <div>
                <div className="header">Your Addresses</div>
                <div className="subtext">
                  Edit, delete, and create shipping addresses
                </div>
              </div>
            </div>
            <div className="dashboard-grid-item"
              onClick={() => navigate("/dashboard/furniture-listings")}
            >
              <HouseIcon />
              <div>
                <div className="header">Your Furniture Listings</div>
                <div className="subtext">
                  View all of your furniture listings
                </div>
              </div>
            </div>
            {/* <div>UserID: {userInfo?.UserID}</div>
            <div>Username: {userInfo?.username}</div>
            <div>Email: {userInfo?.email}</div>
            <div>Password: {userInfo?.password}</div>
            <div>Phone: {userInfo?.phone}</div>
            <div>SessionID: {userInfo?.sessionId}</div>
            <div>Balance: {userInfo?.balance}</div> */}
          </div>
        </div>
      </main>
    </>
  )
}