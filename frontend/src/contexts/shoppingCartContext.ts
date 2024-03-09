import { createContext, useContext } from "react";
import { FurnitureListing } from "../pages/Market";



export type Cart = { [key: string]: FurnitureListing }

export type ShoppingCartState = {
    cart: Cart
    setCart: React.Dispatch<React.SetStateAction<Cart>>
}

export const ShoppingCartContext = createContext<ShoppingCartState | undefined>(undefined)

export function useShoppingCartContext() {
    const cartState = useContext(ShoppingCartContext)

    if (cartState === undefined) {
        throw new Error("useShoppingCartContext is missing")
    }

    return cartState
}