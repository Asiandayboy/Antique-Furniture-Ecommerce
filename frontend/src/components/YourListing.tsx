import { FurnitureListing } from "../pages/Market";
import { convertBase64ToImage } from "../util/image";

export default function YourListing(listing: FurnitureListing) {
  const imageURL = convertBase64ToImage(listing.images[0])

  return (
    <div className="your-listing_wrapper">
      <div className="your-listing-main">
        <div>
          <div className="img_wrapper">
            <img src={imageURL} alt="furniture listing image 1" />
          </div>
          <div className="listing-header">
            <h3 className="listing-title">{listing.title}</h3>
            <div className="listing-desc">{listing.description}</div>
          </div>
        </div>
        {
          listing.bought && <div className="listing-bought">SOLD</div> 
          ||
          <div className="listing-cost">${listing.cost}</div>
        }
      </div>
      <div className="listing-metadata">
        <div className="listing-type">Type: {listing.type}</div>
        <div className="listing-material">Material: {listing.material}</div>
        <div className="listing-condition">Condition: {listing.condition}</div>
        <div className="listing-style">Style: {listing.style}</div>
      </div>
    </div>
  )
}