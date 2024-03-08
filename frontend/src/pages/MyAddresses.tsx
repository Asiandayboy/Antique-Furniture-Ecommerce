import { useEffect, useState } from "react";
import Navbar from "../components/Navbar";



export type ShippingAddress = {
  addressId?: string,
  userId?: string,
  state: string,
  city: string,
  street: string,
  zipCode: string,
  default: boolean
}




export default function MyAddresses() {
  const [addresses, setAddresses] = useState<ShippingAddress[]>([])
  const [editMode, setEditMode] = useState<boolean>(false)
  const [formInput, setFormInput] = useState<ShippingAddress>({
    street: "",
    city: "",
    state: "",
    zipCode: "",
    default: false
  })

  useEffect(() => {
    if (editMode == false) {
      fetch("http://localhost:3000/account/address", {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include"
      })
      .then((res: Response) => res.json())
      .then((data: ShippingAddress[]) => {
        console.log("returned data:", data)
        if (data) {
          setAddresses(data)
        }
      })
      .catch((err: Error) => {
        console.error(err)
      })
    }
  }, [editMode])

  function addAddress(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    console.log("form:", formInput)

    fetch("http://localhost:3000/account/address", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify(formInput)
    })
    .then(async (res: Response) => {
      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      }
      return res.text()
    })
    .then((msg: String) => {
      console.log(msg)
    })
    .catch((err: Error) => {
      console.error(err)
    })


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
              <div>
                <div>ID: {address.addressId}</div>
                <div>Street: {address.street}</div>
                <div>City: {address.city}</div>
                <div>State: {address.state}</div>
                <div>ZipCode: {address.zipCode}</div>
                <div>Default: {String(address.default)}</div>
              </div>
            ))}
          </div> ||
          <div className="address-form_wrapper">
            <form onSubmit={addAddress} className="address-form">
              <h1>Add a new shipping address</h1>
              <div className="first-fields">
                <label htmlFor="">Address</label>
                <input onChange={
                  (e) => setFormInput({...formInput, street: e.currentTarget.value})
                } className="street" type="text" placeholder="Street address" />
              </div>
              <div className="second-fields">
                <input onChange={
                  (e) => setFormInput({...formInput, city: e.currentTarget.value})
                }  className="city" type="text" placeholder="City" />
                <input onChange={
                  (e) => setFormInput({...formInput, state: e.currentTarget.value})
                }  className="state" type="text" placeholder="State" />
                <input onChange={
                  (e) => setFormInput({...formInput, zipCode: e.currentTarget.value})
                }  className="zip" type="text" placeholder="Zip" />
              </div>
              <div>
                <input onChange={
                  (e) => setFormInput({...formInput, default: Boolean(e.currentTarget.value)})
                }  type="checkbox" />
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