import Navbar from "../components/Navbar"
import { Link } from "react-router-dom"

type Props = {
  isLoggedIn: boolean
}


export default function Home({ isLoggedIn }: Props) {
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
          <div className="hero-latest">
            <div className="hero-latest_header">Latest Furniture Listed</div>
            <div className="hero-img">

            </div>
          </div>
        </section>
      </main>
    </>
  )
}