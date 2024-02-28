import Balance from "./Balance.js";
import WalletTransactions from "./WalletTransactions.js";
import AddMoney from "./AddMoney.js";

import * as api from '../Api.js'


const Wallet = () => {


  return (
    <div>
      <Balance></Balance>
      <WalletTransactions></WalletTransactions>
      <AddMoney></AddMoney>
    </div>
  );
};

export default Wallet;
