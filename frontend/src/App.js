import React, { useEffect, useState } from "react";
import FaucetPage from "./IncoFaucet";

function App() {
    const [chains, setChains] = useState([]);
    const [error, setError] = useState(null);
    const [alert, setAlert] = useState(null);
    const [success, setSuccess] = useState(null);
    const [isSending, setIsSending] = useState(false);
    const [address, setAddress] = useState("");
    const [selected, setSelected] = useState(null);


    const loadChains = async () => {
      return fetch("/api/v1/faucet/chains").then((response) => response.json());
    };

    const handleSubmit = async (event) => {
      event.preventDefault();
      setIsSending(true);
      setSuccess(null);
      setAlert(null);
      let payload = {
        address: address,
        chainId: selected.chainId,
      };
      const response = await fetch("/api/v1/faucet/give_me", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });
      let res_json = await response.json();

      setIsSending(false);
      response.status === 200
        ? setSuccess("your request is sent, it may takes some time to proceed")
        : setAlert(res_json.message);
      response.status === 200 ? setAlert(null) : setSuccess(null);
      console.log(response);
    };

    useEffect(() => {
      loadChains()
        .then((res) => {
          setChains(res.chains);
          setSelected(res.chains[0]);
        })
        .catch((err) => {
          console.error(err);
          setError("server is not available");
        });
    }, []);

  return (
    <div>

      <FaucetPage
        isSending={isSending}
        address={address}
        onAddressChanged={setAddress}
        handleSubmit={handleSubmit}
        alert={alert}
        error={error}
        success={success}
        chains={chains}
      />
    </div>
  );
}

export default App;
