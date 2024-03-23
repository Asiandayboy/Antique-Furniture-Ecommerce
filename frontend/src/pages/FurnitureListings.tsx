import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";
import { FurnitureListing } from "./Market";

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
      <div>
        <h1>Your Furniture Listings</h1>
        {listings.map((listing) => (
          <div key={listing.listingID}>
            <div>Title: {listing?.title}</div>
            <div>Description: {listing?.description}</div>
            <div>Cost: {listing?.cost}</div>
            <div>Material: {listing?.material}</div>
            <div>Style: {listing?.style}</div>
            <div>Type: {listing?.type}</div>
            <div>Condition: {listing?.condition}</div>
            <div>Bought: {String(listing?.bought)}</div>
            <div>ListingID: {listing?.listingID}</div>
            <div>SellerID: {listing?.userID}</div>
          </div>
        ))}
      </div>
    </>
  )
}