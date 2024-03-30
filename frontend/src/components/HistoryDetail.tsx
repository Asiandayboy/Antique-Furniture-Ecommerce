import { FurnitureListing } from "../pages/Market";
import { convertBase64ToImage } from "../util/image";
import ImageSlider from "./ImageSlider";


type Props = {
  data: FurnitureListing
  datePurchased: string
  itemNumber: number
}


export default function HistoryDetail({ data, datePurchased, itemNumber }: Props) {

  const iamgeURLs = data.images.map((url) => convertBase64ToImage(url))

  return (
    <div className="history-detail_wrapper">
      <div className="detail-tag">Item #{itemNumber}</div>
      <div className="detail-header">
        <h2>{data.title}</h2>
        <div className="desc-content">Description: <span>{data.description}</span></div>
      </div>
      <div className="detail-item">
        <div>
          <ImageSlider imageURLs={iamgeURLs}/>
          <div className="detail-main">
            <div className="detail-orderid">OrderID: #{data.listingID}</div>
            <div className="detail-order-info">
              <div>
                <div className="info-label">SellerID:</div>
                <div className="info-content">{data.userID}</div>
              </div>
              <div>
                <div className="info-label">Price:</div>
                <div className="info-content">${data.cost}</div>
              </div>
              <div>
                <div className="info-label">Date:</div>
                <div className="info-content">{datePurchased}</div>
              </div>
            </div>
            <div className="detail-metadata">
              <div className="listing-type">Type: <div>{data.type}</div></div>
              <div className="listing-material">Material: <div>{data.material}</div></div>
              <div className="listing-condition">Condition: <div>{data.condition}</div></div>
              <div className="listing-style">Style: <div>{data.style}</div></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}