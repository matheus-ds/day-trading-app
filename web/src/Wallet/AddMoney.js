import React, { useState, useEffect } from 'react';
import { InputGroup, Form, Button } from 'react-bootstrap';
import { useNavigate } from "react-router-dom";
import * as api from '../Api.js'

function AddMoney() {
    const [loading, setLoading] = useState(false);
    const [val, setVal] = useState(0);
    const navigate = useNavigate();

    async function submit() {
        setVal('')
        setLoading(true)
            let msg = await api.addMoneyToWallet(val);
            console.log(msg.success)
            if (msg.success) {
                navigate(0);
            } else {
                alert(msg.message);
            }
      
        setLoading(false)
    }

  return (
    <div style={{marginTop: 70}}>
      <h1>Add Money</h1>
        {loading ? (<h3>Loading....</h3>) : (
      <InputGroup className="mb-3">
        <Button onClick={submit} variant="outline-secondary" id="button-addon1">
          Submit
        </Button>
        <Form.Control onChange={(e) => setVal(e.target.value)}
        placeholder="Add money to wallet"
          aria-label="Example text with button addon"
          aria-describedby="basic-addon1"
        />
      </InputGroup>)}
    </div>
  );
}

export default AddMoney