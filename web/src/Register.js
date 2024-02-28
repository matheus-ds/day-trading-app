import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useNavigate } from "react-router-dom";

const Register = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = (event) => {
    event.preventDefault();
    // Here you can add your Register logic
    console.log('Email:', email);
    console.log('Password:', password);
    // Reset the form
    setEmail('');
    setPassword('');
  };

  return (
    <div className="container">
      <div className="row justify-content-center">
        <div className="col-md-6">
          <h2 className="text-center mb-4">Register</h2>
          <Form onSubmit={handleSubmit}>
            <Form.Group controlId="formBasicEmail">
              <Form.Label>Email address</Form.Label>
              <Form.Control
                type="email"
                placeholder="Enter email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
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
            <Button variant="primary" type="submit" block>
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

