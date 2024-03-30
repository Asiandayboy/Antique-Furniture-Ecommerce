import { useEffect, useState } from "react"
import Navbar from "../components/Navbar"
import { Link } from "react-router-dom"
import { FurnitureListing } from "./Market"
import LatestListing from "../components/LatestListing"

type Props = {
  isLoggedIn: boolean
}


export default function Home({ isLoggedIn }: Props) {
  const [mostRecentListing, setMostRecentListing] = useState<FurnitureListing | null>(null)



  useEffect(() => {

    const fetchMostRecentListing = async () => {
      try {
        const res = await fetch("http://localhost:3000/recent_listing", {
          method: "GET",
          headers: {
            "Content-Type": "application/json"
          },
        })

        if (!res.ok) {
          const msg = await res.text()
          throw new Error(msg)
        }

        const data = await res.json()
        setMostRecentListing(data)
        console.log("recent listing:", data)

      } catch(err) {
        console.error(err)
      }
    }

    fetchMostRecentListing()


  }, [])

  return (
    <>
      <Navbar />
      <main>
        <section className="hero-wrapper">
          <div className="hero-info">
            <div>
              <h1>Shop for great Antique Furnitures!</h1>
              <Link to="/market">
                <button className="begin-shopping_btn">Begin Shopping</button>
              </Link>
            </div>
          </div>
          {
            mostRecentListing &&
            <LatestListing listing={mostRecentListing}/>
          }
        </section>
      </main>
    </>
  )
}