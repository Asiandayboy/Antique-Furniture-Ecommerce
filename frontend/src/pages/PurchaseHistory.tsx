import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";
import { useNavigate } from "react-router-dom";
import { ShippingAddress } from "./MyAddresses";



type OrderItem = {
  orderId: string,
  shippingAddress: ShippingAddress[],
  paymentMethod: string,
  totalCost: number,
  items: string[],
  userId: string,
  datePurchased: string,
  estimatedDelivery: string,
}



export default function PurchaseHistory() {
  const [orders, setOrders] = useState<OrderItem[]>([])


  const navigate = useNavigate()

  function onItemClick(e: React.MouseEvent<HTMLDivElement>) {
    navigate("/dashboard/purchase-history/666")
  }

  useEffect(() => {
    fetch("http://localhost:3000/account/purchase_history", {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include"
    })
    .then(async (res: Response) => {
      if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg)
      } else {
        return res.json()
      }
    })
    .then((data: OrderItem[]) => {
      if (data == null) {
        console.log("No data yet")
      } else {
        setOrders(data)
        console.log("Data:", data)
      }
    })
    .catch((err: Error) => {
      console.error(err)
    })



  }, [])


  return (
    <>
      <Navbar />
      <main>
        <div className="purchase-history">
          <div className="purchase-history-labels">
            <div>Date</div>
            <div>Title</div>
            <div>Amount</div>
          </div>
          <div className="purchase-history_wrapper">
            {
              orders.length > 0 &&
              orders.map((order) => (
                <div onClick={(e) => {
                  navigate(`/dashboard/purchase-history/${order.orderId}`)
                }} className="purchase-history_item">
                  <div className="date">{order.datePurchased}</div>
                  <div className="title">{order.orderId}</div>
                  <div className="amount">{order.totalCost}</div>
                </div>
              )) ||
              <div className="purchase-history_item-null">
                <div className="no-history">
                  You have not made any orders yet
                </div>
              </div>
            }
          </div>
        </div>
      </main>
    </>
  )
}