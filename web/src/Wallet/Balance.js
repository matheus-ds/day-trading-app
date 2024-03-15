import React, { useState, useEffect } from 'react';
import * as api from '../Api.js'


const Balance = () => {
  const [err, setErr] = useState(false);
  const [val, setVal] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {

    // Simulating an API call with setTimeout
    const fetchData = async () => {
      try {
        //await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Fetch data from API
        const data = await api.getWalletBalance();
        setErr(!data.success);
        if (data.success) {
            setVal(data.data.balance ?? []);
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

  return (
    <div>
      {loading ? (
        <h1>{'Balance: Loading....'}</h1>
      ) : (
        <div>
          <h1>{'Balance: '  + (err ? '(error) ' : '') + val}</h1>
        </div>
      )}
    </div>
  );
};

export default Balance;