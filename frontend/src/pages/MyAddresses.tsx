import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";



type ShippingAddress = {
  id: string,
  userId: string,
  state: string,
  city: string,
  street: string,
  zipCode: string,
  default: string
}




export default function MyAddresses() {
  const [addresses, setAddresses] = useState<ShippingAddress[]>([])
  const [editMode, setEditMode] = useState<boolean>(false)

  useEffect(() => {
    fetch("http://localhost:3000/account/address", {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include"
    })
    .then((res: any) => res.json())
    .then((data: ShippingAddress[]) => {
      console.log(data)
      if (data) {
        setAddresses(data)
      }
    })
  })

  function addAddress(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    setEditMode(false)
  }

  return (
    <>
      <Navbar />
      <main>
        {
          !editMode &&
          <div className="address_wrapper">
            {
              addresses.length < 9 &&
              <button 
                onClick={(e) => setEditMode(true)} 
                className="add-addr_btn">+ Add Address</button>
            }
            {addresses.map((address) => (
              <div>{address.id}</div>
            ))}
            {/* <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div>
            <div>bruh</div> */}
          </div> ||
          <div className="address-form_wrapper">
            <form onSubmit={addAddress} className="address-form">
              <h1>Add a new shipping address</h1>
              <div className="first-fields">
                <label htmlFor="">Address</label>
                <input className="street" type="text" placeholder="Street address" />
              </div>
              <div className="second-fields">
                <input className="city" type="text" placeholder="City" />
                <input className="state" type="text" placeholder="State" />
                <input className="zip" type="text" placeholder="Zip" />
              </div>
              <div>
                <input type="checkbox" />
                <label htmlFor="">Make this my default address</label>
              </div>
              <div className="button_wrapper">
                <button type="submit">Add address</button>
                <button onClick={(e) => setEditMode(false)}>Cancel</button>
              </div>
            </form>
          </div>
        }
      </main>
    </>
  )
}