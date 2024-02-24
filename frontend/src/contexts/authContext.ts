import { createContext, useContext } from "react";


export type AuthState = {
    auth: boolean,
    setAuth: React.Dispatch<React.SetStateAction<boolean>>
}

export const AuthContext = createContext<AuthState | undefined>(undefined)

export function useAuthContext() {
    const authState = useContext(AuthContext)

    if (authState === undefined) {
        throw new Error("useAuthContext is missing")
    }

    return authState
}