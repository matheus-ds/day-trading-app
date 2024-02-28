
import styled from 'styled-components'
import Login from './Login'
import {Button, Nav, Navbar, NavDropdown} from 'react-bootstrap';  
import { Outlet, Link, useNavigate } from 'react-router-dom';

function Main() {
  const navigate = useNavigate();
  return (  
    <div>
      <Navbar bg="light" expand="lg">
      <Navbar.Brand href="#home">Day trading app</Navbar.Brand>
      <Navbar.Toggle aria-controls="basic-navbar-nav" />
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto">
          <Nav.Link onClick={() => navigate('/')} >Wallet</Nav.Link>
          <Nav.Link onClick={() => navigate('/mystock')}>Stock portfolio</Nav.Link>
          <Nav.Link onClick={() => navigate('/stockorder')}>Stock order</Nav.Link>
          <Nav.Link onClick={() => navigate('/cancelstock')}>Cancel stock</Nav.Link>
        </Nav>
      </Navbar.Collapse>
      </Navbar> 
      <Outlet></Outlet>
    </div>
  );  
}

export default Main