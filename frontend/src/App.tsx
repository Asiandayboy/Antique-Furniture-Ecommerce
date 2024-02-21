import { useState } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css'
import Home from './pages/Home';
import Signup from './pages/Signup';
import Login from './pages/Login';
import Market from './pages/Market';
import SignupSuccess from './pages/SignupSuccess';

function App() {

  return (
    <BrowserRouter>
      <Routes>
        <Route index element={<Home />}/>
        <Route path='/signup' element={<Signup />} />
        <Route path='/signup-success' element={<SignupSuccess />} />
        <Route path='/login' element={<Login />} />
        <Route path='/market' element={<Market />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
