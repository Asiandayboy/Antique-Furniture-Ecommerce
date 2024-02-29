import { useEffect, useState } from "react"
import MarketFilter from "../components/MarketFilter"
import Navbar from "../components/Navbar"
import FurnitureListing from "../components/FurnitureListing"

type Props = {
  isLoggedIn: boolean
}


type FurnitureListing = {
  listingID: string,
  title: string,
  description: string,
  cost: string,
  type: string,
  style: string,
  condition: string,
  material: string,
  images: string[],
  userID: string,
  bought: string,
}


export default function Market({ isLoggedIn }: Props) {
  const [data, setData] = useState<FurnitureListing[]>([
    // {
    //   listingID: "4444",
    //   title: "bruh",
    //   description: "hello world",
    //   cost: "44..44",
    //   type: "bed",
    //   style: "sheraton",
    //   condition: "mint",
    //   material: "walnut",
    //   images: ["iamge1", "image2"],
    //   userID: "41231231",
    //   bought: "false",
    // }
  ])

  useEffect(() => {
    console.log("FETCHING FURNITURE LISTINGS");

    fetch("http://localhost:3000/get_furnitures", {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      },
    })
      .then(res => res.json())
      .then((data: FurnitureListing[]) => {
        console.log(data);
        setData(data)
      })
      .catch((err: any) => {
        console.log("Error caught:", err)
      })





  }, [])
  
  return (
    <>
      <Navbar />
      <main>
        <section className="market_wrapper">
          <MarketFilter />
          <div className="furniture-listings_wrapper">
            {
              data.map((listing) => (
                <FurnitureListing key={listing.listingID} {...listing} />
              ))
            }
          </div>
        </section>
      </main>
    </>
  )
}