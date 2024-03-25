import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";
import { FurnitureListing } from "./Market";
import YourListing from "../components/YourListing";

export default function FurnitureListings() {
  const [listings, setListings] = useState<FurnitureListing[]>([])


  async function fetchFurnitureListingsOfUser() {
    try {
      const res = await fetch("http://localhost:3000/account/furniture_listings", {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include"
      })

      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      }

      const data = await res.json()
      setListings(data)
      console.log("User's furniture listings:", data)

    } catch (err) {
      console.error(err)
    }
  }



  useEffect(() => {
    fetchFurnitureListingsOfUser()
  }, [])

  return (
    <>
      <Navbar />
      <div className="your-furniture-listings_wrapper">
        <h1>Your Furniture Listings</h1>
        <div className="listings_wrapper">
          {listings.map((listing) => (
            <YourListing {...listing} key={listing.listingID}/>
          ))}
        </div>
      </div>
    </>
  )
}