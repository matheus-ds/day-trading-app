import React, { useState, useEffect } from 'react';
import { InputGroup, Form, Button } from 'react-bootstrap';
import * as api from '../Api.js'

function CancelStock() {
    const [loading, setLoading] = useState(false);
    const [val, setVal] = useState('');

    async function submit() {
        setVal('')
        setLoading(true)
            let msg = await api.cancelStock(val);
            console.log(msg.success)
            if (msg.success == false) {
                alert(msg.data.error);
            }
      
        setLoading(false)
    }

  return (
    <div>
        <h2 className="text-center mb-4">Cancel stock order</h2>

        {loading ? (<h3>Loading....</h3>) : (
      <InputGroup className="mb-3">
        <Button onClick={submit} variant="outline-secondary" id="button-addon1">
          Submit
        </Button>
        <Form.Control onChange={(e) => setVal(e.target.value)}
        placeholder="Stock id"
          aria-label="Example text with button addon"
          aria-describedby="basic-addon1"
        />
      </InputGroup>)}
    </div>
  );
}

export default CancelStock