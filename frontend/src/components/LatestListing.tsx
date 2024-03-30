import { useState } from "react"
import { FurnitureListing } from "../pages/Market"
import { convertBase64ToImage } from "../util/image"
import ImageSlider from "./ImageSlider"
import { Link } from "react-router-dom"

type Props = {
  listing: FurnitureListing
}

export default function LatestListing({ listing }: Props) {
  // const [imageURL, setImageURL] = useState<string>(convertBase64ToImage(listing.images[0]))

  return (
    <div className="latest-listing_wrapper">
      <h2>Latest Furniture Listed</h2>
      <Link to={`/market/${listing.listingID}`}>Go to listing</Link>
      <ImageSlider imageURLs={listing.images.map((url) => convertBase64ToImage(url))}/>
      {/* <img src={imageURL} alt="image of most recent furniture listing" /> */}
    </div>
  )
}