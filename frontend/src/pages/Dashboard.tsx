import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import Navbar from "../components/Navbar"
import HouseIcon from "../assets/HouseIcon"
import { useAccountDataContext } from "../contexts/accountDataContext"

type Props = {
  // isLoggedIn: boolean
}



export default function Dashboard({  }: Props) {

  const navigate = useNavigate()
  const { userData } = useAccountDataContext()

  return (
    <>
      <Navbar />
      <main>
        <div className="dashboard_wrapper">
          <h1>Hello, {userData?.username}! Welcome to you Dashboard</h1>
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
          </div>
        </div>
      </main>
    </>
  )
}