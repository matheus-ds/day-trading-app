import React, { useState, useEffect } from 'react';
import { Table } from 'react-bootstrap';

import * as api from '../Api.js'


const WalletTransactions = () => {
  const [err, setErr] = useState(false);
  const [val, setVal] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {

    // Simulating an API call with setTimeout
    const fetchData = async () => {
      try {
        //await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Fetch data from API
        const data = await api.getWalletTransactions();
        setErr(!data.success);
        if (data.success) {
          setVal(data.data);
        } else {
            setVal(data.message);
        }
        setLoading(false); 
      } catch (error) {
        console.error('Error fetching data:', error);
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  let rows = ["wallet_tx_id","stock_tx_id","is_debit","amount","time_stamp"]


  return (
    <div style={{marginTop: 40}}>
      <h1>{'Transactions'}</h1>
      {loading ? (
        <h3>{'Loading....'}</h3>
      ) : (err ? (<h2>'(error) '{val}</h2>) : (
        <div>
        <h2>Transaction Data</h2>
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
    </div>)
      )
    }
    </div>
  );
};

export default WalletTransactions;