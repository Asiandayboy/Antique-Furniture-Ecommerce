import { createContext, useContext } from "react";


export type AccountInfo = {
    UserID: string,
    username: string,
    email: string,
    password: string,
    phone: string,
    sessionId: string,
    balance: string,
    subscribed: boolean
}


export type AccountData = {
    userData: AccountInfo | null,
    setUserData: React.Dispatch<React.SetStateAction<AccountInfo | null>>
}


export const AccountDataContext = createContext<AccountData | undefined>(undefined)

export function useAccountDataContext() {
    const accountData = useContext(AccountDataContext)

    if (accountData === undefined) {
        throw new Error("useAccountDataContext is missing")
    }

    return accountData
}