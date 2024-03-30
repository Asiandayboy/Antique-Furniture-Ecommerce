import { useShoppingCartContext } from "../contexts/shoppingCartContext"
import { FurnitureListing } from "../pages/Market"
import { convertBase64ToImage } from "../util/image"

type Props = FurnitureListing

export default function CartItem(item: Props) {

  const { cart, setCart }  = useShoppingCartContext()

  const imgURL = convertBase64ToImage(item.images[0])

  function removeFromCart() {
    const updatedCart = { ...cart }
    delete updatedCart[item.listingID]

    setCart(updatedCart)
  }


  return (
    <div className="cart-item">
      <div className="cart-main">
        <div className="cart-item-img_wrapper"><img src={imgURL} alt={`Image of ${item.title}`} /></div>
        <div className="cart-header">
          <h3>{item.title}</h3>
          <div>{item.description}</div>
        </div>
      </div>
      <div className="cart-side">
        <div>Cost: ${item.cost}</div>
        <button onClick={removeFromCart}>Remove from cart</button>
      </div>
    </div>
  )
}