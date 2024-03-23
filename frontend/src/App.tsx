import { useEffect, useState, createContext } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css'
import Home from './pages/Home';
import Signup from './pages/Signup';
import Login from './pages/Login';
import Market from './pages/Market';
import SignupSuccess from './pages/SignupSuccess';
import Dashboard from './pages/Dashboard';
import ProtectedRoutes from './components/ProtectedRoutes';
import { AuthContext } from './contexts/authContext';
import PurchaseHistory from './pages/PurchaseHistory';
import FurnitureListings from './pages/FurnitureListings';
import MyAddresses from './pages/MyAddresses';
import ListFurniture from './pages/ListFurniture';
import PurchaseHistoryDetails from './pages/PurchaseHistoryDetails';
import DetailedListing from './pages/DetailedListing';
import { Cart, ShoppingCartContext } from './contexts/shoppingCartContext';
import ShoppingCart from './pages/ShoppingCart';
import CheckoutSuccess from './pages/checkout/CheckoutSuccess';
import CheckoutCanceled from './pages/checkout/CheckoutCanceled';
import LoginSecurity from './pages/LoginSecurity';
import { AccountInfo, AccountDataContext } from './contexts/accountDataContext';



function App() {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)
  const [cart, setCart] = useState<Cart>({})
  const [userData, setUserData] = useState<AccountInfo | null>(null)

  useEffect(() => {
    const loggedIn = localStorage.getItem("isLoggedIn")

    if (isLoggedIn) {
      localStorage.setItem("isLoggedIn", "true")
    } else if (!isLoggedIn && loggedIn == "true") {
      setIsLoggedIn(true)
    } else {
      localStorage.setItem("isLoggedIn", "false")
      setCart({})
    }
  }, [isLoggedIn])


  return (
    <AuthContext.Provider value={{
      auth: isLoggedIn,
      setAuth: setIsLoggedIn
    }}>
      <AccountDataContext.Provider value={{
        userData: userData,
        setUserData: setUserData
      }}>
        <ShoppingCartContext.Provider value={{
          cart: cart,
          setCart: setCart,
        }}>
          <BrowserRouter>
            <Routes>
              {/* public routes */}
              <Route path="/" element={<Home isLoggedIn={isLoggedIn} />}/>
              <Route path='/market' element={<Market isLoggedIn={isLoggedIn} />} />
              <Route path='/market/listing/:listingId' element={<DetailedListing />}/>
              <Route path='/list' element={<ListFurniture />}/>
              <Route path="/shopping-cart" element={<ShoppingCart />}/>
              <Route path="/checkout_success" element={<CheckoutSuccess />}/>
              <Route path="/checkout_cancel" element={<CheckoutCanceled />}/>

              {/* Routes only accessible when logged out */}
              <Route element={<ProtectedRoutes auth={!isLoggedIn} redirect='/dashboard' />}>
                <Route path='/signup' element={<Signup />} />
                <Route path='/signup-success' element={<SignupSuccess />} />
                <Route path='/login' element={<Login setIsLoggedIn={setIsLoggedIn} />} />
              </Route>

              {/* Routes only accessible when logged in */}
              <Route element={<ProtectedRoutes auth={isLoggedIn} redirect='/login' />}>
                <Route path='/dashboard' element={<Dashboard />} />
                <Route path='/dashboard/purchase-history' element={<PurchaseHistory />} />
                <Route path='/dashboard/purchase-history/:orderId' element={<PurchaseHistoryDetails />} />
                <Route path='/dashboard/furniture-listings' element={<FurnitureListings />} />
                <Route path='/dashboard/addresses' element={<MyAddresses />} />
                <Route path='dashboard/login-security' element={<LoginSecurity />}/>
              </Route>
            </Routes>
          </BrowserRouter>
        </ShoppingCartContext.Provider>
      </AccountDataContext.Provider>
    </AuthContext.Provider>
  )
}

export default App
