import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useNavigate } from "react-router-dom";
import * as api from './Api.js'


const Register = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (event) => {
    event.preventDefault();
    // Here you can add your Register logic
    console.log('Email:', email);
    console.log('Password:', password);
    // Reset the form
    setEmail('');
    setName('');
    setPassword('');
  };

  async function register(email, name, password) {
    let p = await api.register(email, name, password)
    alert("success : " + p.success)
  }

  return (
    <div className="container">
      <div className="row justify-content-center">
        <div className="col-md-6">
          <h2 className="text-center mb-4">Register</h2>
          <Form onSubmit={handleSubmit}>
            <Form.Group controlId="formBasicEmail">
              <Form.Label>User name</Form.Label>
              <Form.Control
                type="user name"
                placeholder="user name"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </Form.Group>
            <Form.Group>
              <Form.Label>Name</Form.Label>
              <Form.Control
                type="name"
                placeholder="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </Form.Group>
            <Form.Group controlId="formBasicPassword">
              <Form.Label>Password</Form.Label>
              <Form.Control
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Form.Group>

            <div style={{paddingTop: 30}}>
            <Button onClick={() => register(email, name, password)} variant="primary" type="submit" block>
              Register
            </Button>
            <Button onClick={() => navigate('/login')}
             style={{marginLeft: 10}} variant="primary" type="submit" block>
              Already resigterd? Sign in
            </Button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
};

export default Register;

