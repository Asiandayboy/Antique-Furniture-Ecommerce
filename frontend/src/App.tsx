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



function App() {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)

  useEffect(() => {
    const loggedIn = localStorage.getItem("isLoggedIn")
    console.log("localstorage:", loggedIn)

    if (loggedIn == "true") {
      setIsLoggedIn(true)
    } else {
      setIsLoggedIn(false)
    }
  }, [])


  return (
    <AuthContext.Provider value={{
      auth: isLoggedIn,
      setAuth: setIsLoggedIn
    }}>
      <BrowserRouter>
        <Routes>
          {/* public routes */}
          <Route path="/" element={<Home isLoggedIn={isLoggedIn} />}/>
          <Route path='/market' element={<Market isLoggedIn={isLoggedIn} />} />

          {/* Routes only accessible when logged out */}
          <Route element={<ProtectedRoutes auth={!isLoggedIn} redirect='/dashboard' />}>
            <Route path='/signup' element={<Signup />} />
            <Route path='/signup-success' element={<SignupSuccess />} />
            <Route path='/login' element={<Login setIsLoggedIn={setIsLoggedIn} />} />
          </Route>

          {/* Routes only accessible when logged in */}
          <Route element={<ProtectedRoutes auth={isLoggedIn} redirect='/login' />}>
            <Route path='/dashboard' element={<Dashboard />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthContext.Provider>
  )
}

export default App
