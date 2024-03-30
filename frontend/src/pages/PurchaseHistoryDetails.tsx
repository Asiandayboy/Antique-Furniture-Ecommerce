import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar";
import { useEffect, useState } from "react";
import { OrderItem, ProductItem } from "./PurchaseHistory";
import { FurnitureListing } from "./Market";
import HistoryDetail from "../components/HistoryDetail";


/*

Make fetch request to get the FurnitureListing given the ID
from listings collection to display the details of the purchase,
which is the details of the listing itself

*/





export default function PurchaseHistoryDetails() {
  const [order, setOrder] = useState<OrderItem>()
  const [orderItems, setOrderItems] = useState<FurnitureListing[]>([])

  const { orderId } = useParams()

  useEffect(() => {
    async function fetchOrderItems(items: ProductItem[]) {
      try {
        const fetchItemPromises = items.map(async (item): Promise<FurnitureListing> => {
          const res = await fetch(`http://localhost:3000/get_furniture/${item.listingId}`, {
            method: "GET",
            headers: {
              "Content-Type": "application/json"
            },
            credentials: "include"
          })

          if (!res.ok) {
            const msg = await res.text()
            throw new Error(msg);
          }

          return await res.json() as FurnitureListing
        })

        const listings = await Promise.all(fetchItemPromises)

        // discard the failed fetches
        return listings.filter((listing) => listing !== null)
      } catch (err) {
        console.error(err)
        return []
      }
    }

    async function fetchOrder() {
      try {
        const res = await fetch(`http://localhost:3000/account/purchase_history/${orderId}`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json"
          },
          credentials: "include"
        })

        if (!res.ok) {
          const msg = await res.text()
          throw new Error(msg)
        }

        const order: OrderItem = await res.json()
        const orderItems: FurnitureListing[] = await fetchOrderItems(order.items)
        setOrder(order)
        setOrderItems(orderItems)

      } catch (err) {
        console.error(err)
      }
    }

    fetchOrder()


  }, [])

  return (
    <>
      <Navbar />
      <main>
        <div className="purchase-history-detailed_wrapper">
          <h1>Purchase History Item Details</h1>
          <div className="history-list_wrapper">
            {
              orderItems.map((item, i) => (
                <HistoryDetail 
                  data={item} 
                  key={item.listingID}
                  datePurchased={order?.datePurchased!}
                  itemNumber={i+1}
                />
              ))
            }
          </div>
        </div>
      </main>
    </>
  )
}