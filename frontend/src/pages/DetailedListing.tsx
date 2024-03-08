import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar";

export default function DetailedListing() {

  const { listingId } = useParams()


  return (
    <>
      <Navbar />
      <div>
        <div>DetailedListing {listingId} </div>
        <button>Add to cart</button>
      </div>
    </>
  )
}