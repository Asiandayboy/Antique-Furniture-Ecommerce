import { useEffect, useState } from "react"
import MarketFilter from "../components/MarketFilter"
import Navbar from "../components/Navbar"
import FurnitureListing from "../components/FurnitureListing"

type Props = {
  isLoggedIn: boolean
}


export type FurnitureListing = {
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
  const [data, setData] = useState<FurnitureListing[]>([])
  const [filteredData, setFilteredData] = useState<FurnitureListing[]>([])

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
        // console.log(data);
        setData(data)
        setFilteredData(data)
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
          <MarketFilter dataSet={data} setDataSet={setFilteredData} />
          <div className="furniture-listings_wrapper">
            {
              filteredData.map((listing) => (
                <FurnitureListing key={listing.listingID} {...listing} />
              ))
            }
          </div>
        </section>
      </main>
    </>
  )
}

