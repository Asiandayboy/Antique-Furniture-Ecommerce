import { useNavigate } from "react-router-dom"

type Props = {
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


export default function FurnitureListing(data: Props) {
  const navigate = useNavigate()

  function onClick(e: React.MouseEvent<HTMLDivElement>) {
    navigate(`/market/listing/${data.listingID}`)
  }

  return (
    <div onClick={onClick} className="furniture-listing">
      {data.title}
      <div>ListingID: {data.listingID}</div>
      <div>Desc: {data.description}</div>
      <div>Cost: {data.cost}</div>
      <div>Type: {data.type}</div>
      <div>Style: {data.style}</div>
      <div>Material: {data.material}</div>
      <div>Condition: {data.condition}</div>
      <div>Bought: {String(data.bought)}</div>
      <div>SellerID: {data.userID}</div>
    </div>
  )
}