import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { convertBase64ToImage } from "../util/image"

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
  const [imageURLs, setImageURLs] = useState<string[]>([])

  const navigate = useNavigate()

  function onClick(e: React.MouseEvent<HTMLDivElement>) {
    navigate(`/market/${data.listingID}`)
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