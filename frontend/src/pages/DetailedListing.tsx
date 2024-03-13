import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar";
import { useEffect, useState } from "react";
import { FurnitureListing } from "./Market";
import { useShoppingCartContext } from "../contexts/shoppingCartContext";

export default function DetailedListing() {
  const [listingData, setListingData] = useState<FurnitureListing>()
  const [addedToCart, setAddedToCart] = useState<boolean>(false)

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
      console.log("detailed listing:", data)
    })
    .catch((err: Error) => {
      console.error(err)
    })


  }, [])


  useEffect(() => {
    console.log("SHopping cart:", cart)
  }, [cart])


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
      <div>
        <div>
          <div>Title: {listingData?.title}</div>
          <div>Description: {listingData?.description}</div>
          <div>Cost: {listingData?.cost}</div>
          <div>Material: {listingData?.material}</div>
          <div>Style: {listingData?.style}</div>
          <div>Type: {listingData?.type}</div>
          <div>Condition: {listingData?.condition}</div>
          <div>Bought: {String(listingData?.bought)}</div>
          <div>ListingID: {listingData?.listingID}</div>
          <div>SellerID: {listingData?.userID}</div>
        </div>
        {
          !addedToCart &&
          <button className="add-to-cart_btn" onClick={onAddToCart}>Add to cart</button>
          ||
          <div className="post-cart-add_btn">Listing has been added to cart</div>
        }
      </div>
    </>
  )
}