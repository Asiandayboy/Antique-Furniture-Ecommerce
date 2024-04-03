import { AccountInfo, useAccountDataContext } from "../contexts/accountDataContext"
import { getAccountData } from "../util/account"



export default function SubscribeButton() {
  const { userData, setUserData } = useAccountDataContext()

  function onClick() {
    let target = `http://localhost:3000/${userData?.subscribed && "unsubscribe" || "subscribe"}`
    
    fetch(target, {
      method: "POST",
      headers: {
        "Content-Type": "text/plain"
      },
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
      getAccountData()
      .then((val: AccountInfo | Error) => {
        if (val instanceof Error) {
          console.error(val)
        } else {
          setUserData({...val})
        }
      })
      console.log("subscribe result:", data)
    })
    .catch((err) => {
      console.error(err)
    })
  }


  return (
    !userData?.subscribed &&
    <button className="subscribe_btn" onClick={onClick}>
      Subscribe to receive emails of new furniture listings!
    </button> ||
    <button className="subscribe_btn" onClick={(onClick)}>
      Unsubscribe to stop receiving email updates of new furniture listings
    </button>
  )
}