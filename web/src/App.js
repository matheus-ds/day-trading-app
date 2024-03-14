import {
  BrowserRouter,
  Routes,
  Route,
} from 'react-router-dom'
import React, { useState } from 'react';
import { useNavigate, Navigate, Outlet } from "react-router-dom";
import styled from 'styled-components'
import FourOhFour from './404'
import Login from './Login'
import Main from './Main'
import Register from './Register'
import Wallet from './Wallet/Wallet.js'
import MyStock from './MyStock/MyStock.js'
import StockOrder from './StockOrder/StockOrder.js'
import CancelStock from './StockOrder/CancelStock.js'

import * as api from './Api.js'
import { Button } from 'react-bootstrap';



const PageWrapper = styled('div')`
  padding: 24px 70px;
`


function Router() {
  const [authtoken , setAuthtoken ] = useState(false);
  const navigate = useNavigate();

  function authenticate(username, password) {
    api.login(username, password).then(function (response) {
      console.log(response)
      setAuthtoken(response.success)
      if (response.success) {
        api.setToken(response.data.token);
        navigate('/');
      } else {
        alert(response.data.error);
      }
    });
  }

  return (
    <div>

    <Routes>
      <Route element={authtoken ? <Outlet/> : <Navigate to="/login"/>}>
        <Route element={<Main/>}>
          <Route element={<Wallet/>} path="/" exact/>
          <Route element={<MyStock/>} path="/mystock" exact/>
          <Route element={<StockOrder/>} path="/stockorder" exact/>
          <Route element={<CancelStock/>} path="/cancelstock" exact/>

        </Route>
      </Route>
      <Route element={authtoken ? <Navigate to="/"/> : <Outlet/>}>
      <Route element={<Login login={authenticate}/>} path="/login"/>
      <Route path="/register" element={<Register />} />
      </Route>
    </Routes>
    </div>
  )
}



export default function App() {
  return (
      <BrowserRouter>
            <PageWrapper>
              <Router />
            </PageWrapper>
      </BrowserRouter>
  )
}
