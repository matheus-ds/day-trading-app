import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';
import { useNavigate } from "react-router-dom";
import * as api from '../Api.js'

const StockOrder = () => {

  const [stockid, setStockid] = useState('');
  const [ordertype, setOrdertype] = useState(false);
  const [quantity, setQuantity] = useState('');
  const [price, setPrice] = useState('');
  const [loading, setLoading] = useState(false);


  async function handleSubmit (stockid, ordertype, quantity, price) {
    // Reset the form
    setStockid('');
    setQuantity('');
    setPrice('')
    setLoading(true)
    setLoading(true)
    let msg = await api.placeStockOrder(stockid, ordertype, quantity, price);
    console.log(msg.success)
    if (msg.success == false) {
        alert(msg.message);
    }

    setLoading(false)
  };


  return (
    <div>
        {loading ? (<h3>Loading....</h3>) : (
    <div className="container">
    <div className="row justify-content-center">
      <div className="col-md-6">
        <h2 className="text-center mb-4">Place stock order</h2>
        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="stockorder">
            <Form.Label>Stock ID</Form.Label>
            <Form.Control
              placeholder="Stock id"
              value={stockid}
              onChange={(e) => setStockid(e.target.value)}
            />
           <Form.Label>Quantity</Form.Label>
            <Form.Control
              placeholder="quantity"
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
            />
           <Form.Label>Price</Form.Label>
           <Form.Control
              placeholder="price"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
            />
          <Form.Label>Order type</Form.Label>
          <Form.Check // prettier-ignore
              type="checkbox"
              value={true}
              id=""
              label="if checked it is LIMIT else it is MARKET"
              onChange={(e) => {
                setOrdertype(!ordertype)}}
          />
          </Form.Group>

          <div style={{paddingTop: 30}}>
          <Button type="submit" onClick={() => api.placeStockOrder(
            stockid, ordertype ? 'LIMIT' : 'MARKET', quantity, price
          )} variant="primary" block>
            Login
          </Button>
          </div>
        </Form>
      </div>
    </div>
  </div>)
      }
    </div>
  );


};

export default StockOrder;