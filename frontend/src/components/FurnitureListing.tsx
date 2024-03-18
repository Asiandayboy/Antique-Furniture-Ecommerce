import { useEffect, useState } from "react"
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


/**
 * Converts Base64 encoded string of an image 
 * to a BLOB and returns the URL of the image BLOB
 */
function convertBase64ToImage(base64Data: string): string {
  const binaryStr = atob(base64Data)

  const byteArray = new Uint8Array(binaryStr.length)
  for (let i = 0; i < binaryStr.length; i++) {
    byteArray[i] = binaryStr.charCodeAt(i)
  }

  const blob = new Blob([byteArray], { type: "image/png" })
  const url = URL.createObjectURL(blob)

  return url
}


export default function FurnitureListing(data: Props) {
  const [imageURLs, setImageURLs] = useState<string[]>([])

  const navigate = useNavigate()

  function onClick(e: React.MouseEvent<HTMLDivElement>) {
    navigate(`/market/listing/${data.listingID}`)
  }

  useEffect(() => {
    data.images.forEach((imageData) => {
      const imageURL = convertBase64ToImage(imageData)
      setImageURLs([ ...imageURLs, imageURL ])
    })

  }, [])


  return (
    <div onClick={onClick} className="furniture-listing">
      {data.title}
      <div className="img_wrapper">
        <img src={imageURLs[0]} alt="furniture listing image 1" />
      </div>
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