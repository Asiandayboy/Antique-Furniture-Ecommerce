import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar";
import { useEffect, useState } from "react";
import { FurnitureListing } from "./Market";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";
import { convertBase64ToImage } from "../util/image";
import ImageSlider from "../components/ImageSlider";

export default function DetailedListing() {
  const [listingData, setListingData] = useState<FurnitureListing>()
  const [addedToCart, setAddedToCart] = useState<boolean>(false)
  const [imageURLs, setImageURLs] = useState<string[]>([])

  const { listingId } = useParams()
  const { cart, setCart } = useShoppingCartContext()

  // fetch individual furniture listing data from API
  useEffect(() => {
    fetch(`http://localhost:3000/get_furniture/${listingId}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include"
    })
    .then(async (res: Response) => {
      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      } else {
        return res.json()
      }
    })
    .then((data: FurnitureListing) => {
      setListingData(data)
      const urls = data.images.map(imageBase64 => convertBase64ToImage(imageBase64));
      setImageURLs(urls)
      console.log("detailed listing:", data)
    })
    .catch((err: Error) => {
      console.error(err)
    })
  }, [])



  function onAddToCart(e: React.MouseEvent<HTMLButtonElement>) {
    if (!listingId || !listingData) return

    if (!cart[listingId]) {
      setCart(prevCart => ({
        ...prevCart,
        [listingId]: listingData
      }))
      setAddedToCart(true)
    } else {
      alert("This item has already been added to the cart")
    }
  }
  


  return (
    <>
      <Navbar />
      <div className="detailed-listing_wrapper">
        {imageURLs.length > 0 && <ImageSlider imageURLs={imageURLs}/>}
        <div className="detailed-listing-info">
          <div className="detailed-listing-main">
            <h1>{listingData?.title}</h1>
            <div>{listingData?.description}</div>
            {
              listingData?.bought && <div className="detailed-bought">SOLD</div> 
              ||
              <div className="detailed-cost">${listingData?.cost}</div>
            }
          </div>
          <div className="listing-metadata">
            <div className="listing-type">Type: <div>{listingData?.type}</div></div>
            <div className="listing-material">Material: <div>{listingData?.material}</div></div>
            <div className="listing-condition">Condition: <div>{listingData?.condition}</div></div>
            <div className="listing-style">Style: <div>{listingData?.style}</div></div>
          </div>
          {
            !listingData?.bought &&
            (
              !addedToCart &&
              <button className="add-to-cart_btn" onClick={onAddToCart}>Add to cart</button>
              ||
              <div className="post-cart-add_btn">Listing has been added to cart</div>
            )
          }
        </div>
      </div>
    </>
  )
}