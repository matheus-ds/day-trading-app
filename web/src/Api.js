import axios from "axios";
const baseURL = "http://localhost:8000/";
window.token = false

export function hello(){

}

export function setToken(tk){
    window.token = tk
}


export async function login(username, password) {
    let options = {
        method: 'POST',
        url: baseURL + 'login',
        data: {
          user_name: username,
          password: password
        }
    };
    //return await '{"success":true, "data":{"token":"yourtoken"}}');
    var p = (await axios(options)).data
    return p;

}

export async function register(username, name, password) {
    let options = {
        method: 'POST',
        url: baseURL + 'register',
        data: {
            user_name: username,
            password: password,
            name: name
        }
    };
    var pot = await axios(options)
    return pot.data;
}

export async function getWalletBalance() {
    let options = {
        method: 'GET',
        url: baseURL + 'getWalletBalance',
        headers: {
            token : window.token
        }
    };
    //await new Promise(resolve => setTimeout(resolve, 10000));
    // return await 
    //     '{"success":true, "message":"sorry", "data":{"balance": 100}}'
    // );
    var pot = await axios(options)
    return pot.data;
}


export async function getWalletTransactions() {
    console.log(window.token);
    let options = {
        method: 'GET',
        url: baseURL + 'getWalletTransactions',
        headers: {
            token : window.token
        }
    };
    // await new Promise(resolve => setTimeout(resolve, 5000));
    // return await 
    //     `{"success":true, "data":[{"wallet_tx_id":
    //     "628ba23df2210df6c3764823","stock_tx_id":"62738363a50350b1fbb243a6",
    //     "is_debit":true,"amount":100,"time_stamp":"2024-01-12T15:03:25.019+00:00"},
    //     {"wallet_tx_id":"628ba36cf2210df6c3764824",
    //     "stock_tx_id":"62738363a50350b1fbb243a6",
    //     "is_debit":false,"amount":200,
    //     "time_stamp":"2024-01-12T14:13:25.019+00:00"}], "message": "sorry"}`
    // );
    var pot = await axios(options)

    return pot.data;
}

//do not
export async function addMoneyToWallet(amount) {
    let options = {
        method: 'POST',
        url: baseURL + 'addMoneyToWallet',
        headers: {
            token : window.token
        },
        data: {
            amount: amount
        }
    };
    // await new Promise(resolve => setTimeout(resolve, 3000));
    // return await 
    //     '{"success":true, "data":null}'
    // );
    var pot = await axios(options)
    return pot.data;
}

export async function getStockTransactions(amount) {
    let options = {
        method: 'GET',
        url: baseURL + 'getStockTransactions',
        headers: {
            token : window.token
        },
        data: {
            amount: amount
        }
    };
    //await new Promise(resolve => setTimeout(resolve, 3000));
//     return await 
//         `
//         {"success":true, "message":"sorry" , "data":[{
//             "stock_tx_id":"62738363a50350b1fbb243a6",
//         "stock_id":1,"wallet_tx_id":"628ba23df2210df6c764823",
//         "order_status":"COMPLETED","is_buy":true,"order_type":"LIMIT",
// "stock_price":50,"quantity":2,"parent_tx_id": null,
// "time_stamp":"2024-01-12T15:03:25.019+00:00"}]}`
//     );
    var pot = await axios(options)
    return pot.data;
}

export async function getStockPrices(amount) {
    let options = {
        method: 'GET',
        url: baseURL + 'getStockPrices',
        headers: {
            token : window.token
        },
        data: {
            amount: amount
        }
    };
    //await new Promise(resolve => setTimeout(resolve, 3000));
    // return await 
    //     `
    //     {"success":true, "message":"sorry man", "data":[{"stock_id":1,
    //     "stock_name":"Apple","current_price":100},
    //     {"stock_id":1, "stock_name":"Google",
    //     "current_price": 200}]}
    //     `
    // );
    var pot = await axios(options)
    return pot.data;
}

export async function getStockPortfolio(amount) {
    let options = {
        method: 'GET',
        url: baseURL + 'getStockPortfolio',
        headers: {
            token : window.token
        },
        data: {
            amount: amount
        }
    };
    //await new Promise(resolve => setTimeout(resolve, 3000));
    // return await 
    //     `
    //     {"success":true,"message":"nope", "data":[{
    //         "stock_id":1,"stock_name":"Apple",
    //         "quantity_owned":100},{
    //     "stock_id":2,"stock_name":"Google",
    //     "quantity_owned":150}]
    //         }
    // `);
    var pot = await axios(options)
    return pot.data;
}


export async function placeStockOrder(id, type, quantity, price) {
    let options = {
        method: 'POST',
        url: baseURL + 'placeStockOrder',
        headers: {
            token : window.token
        },
        data: {stock_id:id,is_buy:true,order_type: type,
        quantity:quantity,
        price:price}
    };
    // await new Promise(resolve => setTimeout(resolve, 3000));
    // return await 
    //     '{"success":true,"message":"nope", "data":null}'
    // );
    var pot = await axios(options)
    return pot.data;
}

export async function cancelStock(id) {
    let options = {
        method: 'POST',
        url: baseURL + 'cancelStock',
        headers: {
            token : window.token
        },
        data: {
            stock_tx_id: id
        }
    };
    // await new Promise(resolve => setTimeout(resolve, 3000));
    // return await 
    //     '{"success":true, "data":null}'
    // );
    var pot = await axios(options)
    return pot.data;
}