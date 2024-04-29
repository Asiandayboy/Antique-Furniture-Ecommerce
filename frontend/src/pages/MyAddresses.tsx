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

const EMPTY_ADDRESS: ShippingAddress = {
  street: "",
  city: "",
  state: "",
  zipCode: "",
  default: false
}


export default function MyAddresses() {
  const [addresses, setAddresses] = useState<ShippingAddress[]>([])
  const [editMode, setEditMode] = useState<boolean>(false)
  const [formInput, setFormInput] = useState<ShippingAddress>(EMPTY_ADDRESS)
  const [editId, setEditId] = useState<string>("")
  const [key, setKey] = useState<number>(0) // used to force a remount to render updated data
  const [deleteMode, setDeleteMode] = useState<string>("")
  const [errMsg, setErrMsg] = useState<string>("")
  const [isError, setIsError] = useState<boolean>()

  function startErrorAnim() {
    setIsError(true)
    setTimeout(() => {
      setIsError(false)
    }, 1000)
  }

  function isFormValid(input: ShippingAddress): boolean {
    if (!input.street || !input.city || !input.zipCode || !input.street) {
      return false
    }

    return true
  }


  function editAddress(address: ShippingAddress) {
    if (address.addressId && address.addressId.length > 0) {
      setEditMode(true)
      setEditId(address.addressId)

      setFormInput({
        street: address.street,
        city: address.city,
        state: address.state,
        zipCode: address.zipCode,
        default: address.default
      })
    }

  }

  function fetchAddresses() {
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
        setAddresses(data || [])
      })
      .catch((err: Error) => {
        console.error(err)
      })
    }
  }

  function closeForm() {
    setEditMode(false)
    setEditId("")
    setFormInput(EMPTY_ADDRESS)
    setErrMsg("")
    setIsError(false)
  }

  function formSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    if (!isFormValid(formInput)) {
      setErrMsg("Fields cannot be blank")
      startErrorAnim()
      return
    }

    console.log("form:", formInput)

    if (editId.length > 0) {
      fetch("http://localhost:3000/account/address", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({
          addressID: editId,
          changes: formInput
        })
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
        setKey(prev => prev+1)
      })
      .catch((err: Error) => {
        console.error(err)
      })
    } else {
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
        setKey(prev => prev+1)
      })
      .catch((err: Error) => {
        console.error(err)
      })
    }

    closeForm()
  }

  function deleteAddress(addressID: string) {
    fetch(`http://localhost:3000//account/address/${addressID}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
    })
    .then(async (res: Response) => {
      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      }
      return res.text()
    })
    .then((msg: String) => {
      console.log("delete msg:", msg)
      setKey(prev => prev+1)
      setDeleteMode("")
    })
    .catch((err: Error) => {
      console.error(err)
    })

  }


  useEffect(() => {
    fetchAddresses()
  }, [key])

  return (
    <>
      <Navbar />
      <main>
        {
          !editMode &&
          <div className="my-address_wrapper">
            <h1>Your Addresses</h1>
            <div className="addresses">
              {
                addresses.length < 9 &&
                <button 
                  onClick={(e) => setEditMode(true)} 
                  className="add-addr_btn">+ Add Address</button>
              }
              {addresses.map((address) => (
                <div className="address-item" key={address.addressId}>
                  <div>
                    <div>Street: {address.street}</div>
                    <div>City: {address.city}</div>
                    <div>State: {address.state}</div>
                    <div>ZipCode: {address.zipCode}</div>
                    <div>Default: {String(address.default)}</div>
                    <div>ID: {address.addressId}</div>
                  </div>
                  <div className="address-buttons">
                    <button onClick={() => editAddress(address)}>Edit</button>
                    <button onClick={() => setDeleteMode(address.addressId!)}>Delete</button>
                    <dialog 
                      className="address-delete-dialog" 
                      open={deleteMode == address.addressId}
                      style={deleteMode != address.addressId && {"display": "none"} || {}}
                    >
                      <div>
                        <div>Are you sure you want to delete address ID {address.addressId}</div>
                        <div className="dialog-buttons">
                          <button className="yes" onClick={() => deleteAddress(address.addressId!)}>Yes</button>
                          <button className="no" onClick={() => setDeleteMode("")}>No</button>
                        </div>
                      </div>
                    </dialog>
                  </div>
                </div>
              ))}
            </div>
          </div> ||
          <div className="address-form_wrapper">
            <form onSubmit={formSubmit} className="address-form">
              <h1>{editId.length > 0 && "Edit a shipping address" || "Add a new shipping address"}</h1>
              <div className="first-fields">
                <label htmlFor="">Address</label>
                <input onChange={
                  (e) => setFormInput({...formInput, street: e.currentTarget.value})
                } className="street" type="text" placeholder="Street address" value={formInput.street} />
              </div>
              <div className="second-fields">
                <input onChange={
                  (e) => setFormInput({...formInput, city: e.currentTarget.value})
                }  className="city" type="text" placeholder="City" value={formInput.city}/>
                <input onChange={
                  (e) => setFormInput({...formInput, state: e.currentTarget.value})
                }  className="state" type="text" placeholder="State" value={formInput.state}/>
                <input onChange={
                  (e) => setFormInput({...formInput, zipCode: e.currentTarget.value})
                }  className="zip" type="text" placeholder="Zip" value={formInput.zipCode}/>
              </div>
              <div>
                <input onChange={
                  (e) => setFormInput({...formInput, default: Boolean(e.currentTarget.value)})
                }  type="checkbox" checked={formInput.default}/>
                <label htmlFor="">Make this my default address</label>
              </div>
              <div className="button_wrapper">
                <button type="submit">{editId.length > 0 && "Edit address" || "Add address"}</button>
                <button onClick={closeForm}>Cancel</button>
              </div>
              {
                errMsg &&
                <div className={!isError && "addr_err-msg" || "addr_err-msg err-msg-anim"}>{errMsg}</div>
              }
            </form>
          </div>
        }
      </main>
    </>
  )
}