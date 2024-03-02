import React, { useState, useEffect } from 'react';
import { Table } from 'react-bootstrap';

import * as api from '../Api.js'


const StockPortfolio = () => {
  const [err, setErr] = useState(false);
  const [val, setVal] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {

    // Simulating an API call with setTimeout
    const fetchData = async () => {
      try {
        //await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Fetch data from API
        const data = await api.getStockPortfolio();
        setErr(!data.success);
        if (data.success) {
          setVal(data.data);
        } else {
            setVal(data.data.error);
        }
        setLoading(false); 
      } catch (error) {
        setVal(error);
        console.error('Error fetching data:', error);
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  let rows = ["stock id","stock name","quantity owned"]


  return (
    <div style={{paddingTop: 40}}>
      {loading ? (
        <h1>{'Stock portfolio: Loading....'}</h1>
      ) : ( err ? (<h1>{'Stock portfolio (Error): ' + val}</h1>) : (
        <div>
        <h2>Stock portfolio</h2>
        <Table striped bordered hover>
            <thead>
                <tr>
                    {rows.map(key => (
                        <th key={key}>{key}</th>
                    ))}
                </tr>
            </thead>
            <tbody>
            {val.map((item, index) => (
                            <tr key={index}>
                                {Object.values(item).map((value, index) => (
                                    <td key={index}>{value == null ? "null" : value.toString()}</td>
                                ))}
                            </tr>
                        ))}
            </tbody>
        </Table>
    </div>
      ))
    }
    </div>
  );
};

export default StockPortfolio;