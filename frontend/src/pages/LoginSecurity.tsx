import { useState } from "react";
import Navbar from "../components/Navbar";
import { AccountInfo, useAccountDataContext } from "../contexts/accountDataContext";
import { getAccountData } from "../util/account";



type Changes = {
  newPassword?: string,
  newEmail?: string,
  newPhone?: string
}

type EditMode = "None" | "Email" | "Phone" | "Password"



export default function LoginSecurity() {
  const { userData, setUserData } = useAccountDataContext()
  const [changes, setChanges] = useState<Changes>({
    newPassword: "",
    newEmail: userData?.email,
    newPhone: userData?.phone
  })
  const [editMode, setEditMode] = useState<EditMode>("None")


  function saveChanges() {

    let input = {
      newEmail: changes.newEmail,
      newPhone: changes.newPhone,
      newPassword: changes.newPassword
    }


    if (input.newEmail == userData?.email) {
      input = {
        newPhone: input.newPhone,
        newPassword: input.newPassword,
        newEmail: undefined
      }
    }
    if (input.newPhone == userData?.phone) {
      console.log("hone")
      input = {
        newPassword: input.newPassword,
        newPhone: undefined,
        newEmail: input.newEmail
      }
    }
    if (input.newPassword == "") {
      input = {
        newPassword: undefined,
        newPhone: input.newPhone,
        newEmail: input.newEmail
      }
    }



    fetch("http://localhost:3000/account", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(input),
      credentials: "include"
    })
    .then(async (res: Response) => {
      if (!res.ok) {
        const msg = await res.text()
        throw new Error(msg)
      }

      return res.text()
    })
    .then((data: string) => {
      console.log("LoginSec change:", data)

      getAccountData()
      .then((val: AccountInfo | Error) => {
        if (val instanceof Error) {
          console.error(val)
        } else {
          setUserData({...val})
        }
      })
    })
    .catch(err => {
      console.error(err)
    })
    


    setEditMode("None")
  }



  return (
    <>
      <Navbar />
      <div className="login-security_wrapper">
        <h1>LoginSecurity</h1>
        <div className="loginsec-info">
          {
            editMode == "Email" &&
            <div>
              <div>Email: 
                <input type="text" value={changes.newEmail} 
                onChange={(e) => setChanges({...changes, newEmail: e.currentTarget.value.trim()})}/>
              </div>
              <div className="edit_buttons">
                <button onClick={saveChanges}>Save</button>
                <button onClick={() => setEditMode("None")}>Cancel</button>
              </div>
            </div> ||
            <div>
              <div>Email: {userData?.email}</div>
              <button onClick={() => setEditMode("Email")}>Edit</button>
            </div>
          }
          {
            editMode == "Phone" &&
            <div>
              <div>Phone: 
                <input type="tel" value={changes.newPhone} 
                pattern="[0-9]{3}-[0-9]{3}-[0-9]{4}"
                maxLength={12}
                onChange={(e) => setChanges({...changes, newPhone: e.currentTarget.value.trim()})}/>
              </div>
              <div className="edit_buttons">
                <button onClick={saveChanges}>Save</button>
                <button onClick={() => setEditMode("None")}>Cancel</button>
              </div>
            </div> ||
            <div>
              <div>Phone: {userData?.phone == "" ? "No Phone" : userData?.phone}</div>
              <button onClick={() => setEditMode("Phone")}>Edit</button>
            </div>
          }
          {
            editMode == "Password" &&
            <div>
              <div>Password: 
                <input type="password" value={changes.newPassword} 
                onChange={(e) => setChanges({...changes, newPassword: e.currentTarget.value.trim()})}/>
              </div>
              <div className="edit_buttons">
                <button onClick={saveChanges}>Save</button>
                <button onClick={() => setEditMode("None")}>Cancel</button>
              </div>
            </div> ||
            <div>
              <div>Password: ----------</div>
              <button onClick={() => setEditMode("Password")}>Edit</button>
            </div>
          }
        </div>

      </div>
    </>
  )
}