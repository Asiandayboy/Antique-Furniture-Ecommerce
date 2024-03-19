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
      <div className="listing-header">
        <div>
          <div className="listing-title">
            {data.title}
          </div>
          <div className="listing-desc">{data.description}</div>
        </div>
        {
          data.bought && <div className="listing-bought">SOLD</div> 
          ||
          <div className="listing-cost">${data.cost}</div>
        }
      </div>
      <div className="img_wrapper">
        <img src={imageURLs[0]} alt="furniture listing image 1" />
      </div>
      <div className="listing-metadata">
        <div className="listing-type">Type: {data.type}</div>
        <div className="listing-material">Material: {data.material}</div>
        <div className="listing-condition">Condition: {data.condition}</div>
        <div className="listing-style">Style: {data.style}</div>
      </div>
      <span className="overlay-hover">Click to view more</span>
    </div>
  )
}