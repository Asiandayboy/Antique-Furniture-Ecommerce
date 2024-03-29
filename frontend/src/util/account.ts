import { AccountInfo } from "../contexts/accountDataContext"


export function getAccountData(): Promise<AccountInfo | Error> {
    return fetch("http://localhost:3000/account",{
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
        mode: "cors",
        credentials: "include",
    })
    .then(async (res: Response) => {
        if (!res.ok) {
            const msg = await res.text()
            throw new Error(msg)
        }
        return res.json()
    })
    .catch((err: Error) => {
        return err
    })
}