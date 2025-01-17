import React from "react";

const AlertView = ({ message, type }) => {
  const getAlertStyle = () => {
    switch (type) {
      case "error":
        return "bg-red-50 text-red-700 border-red-200";
      case "warning":
        return "bg-yellow-50 text-yellow-700 border-yellow-200";
      case "success":
        return "bg-green-50 text-green-700 border-green-200";
      default:
        return "bg-blue-50 text-blue-700 border-blue-200";
    }
  };

  return (
    <div className={`px-4 py-3 rounded-lg border ${getAlertStyle()} mb-4`}>
      {message}
    </div>
  );
};

function FaucetPage({
  isSending,
  address,
  onAddressChanged,
  handleSubmit,
  error,
  alert,
  success,
  chains,
}) {
  return (
    <div className="flex flex-col md:flex-row h-screen font-urbanist">
      <div className="w-full md:w-1/2 flex flex-col bg-white p-8 md:px-16 md:py-12 h-full">
        <img src="/favicon.ico" className="w-32" alt="Inco Logo" />

        {error ? (
          <div className="max-w-[520px] pt-16 md:pt-32 md:px-6">
            <h1 className="text-[40px] md:text-[4rem] leading-[48px] font-semibold tracking-[-0.03em] text-black mb-6">
              Server Unavailable
            </h1>

            <p className="text-[15px] md:text-xl leading-[22px] font-mono text-gray-900 mb-12">
              We're unable to connect to the server. Please try again later.
            </p>

            <a
              href="https://discord.com/invite/inco"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-[#4477FF] text-sm hover:underline font-mono"
            >
              Visit Help Center
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </a>
          </div>
        ) : (
          <div className="max-w-[520px] pt-16 md:pt-32 md:px-6">
            <h1 className="text-[40px] md:text-[4rem] leading-[48px] font-semibold tracking-[-0.03em] text-black mb-6">
              Faucet
            </h1>

            <p className="text-[15px] md:text-xl leading-[22px] font-mono text-gray-900 mb-12">
              Get some ORAI test tokens.
            </p>

            <p className="text-[13px] leading-[20px] font-normal text-gray-500 font-mono mb-2">
              This is NOT an incentivized testnet, please DO NOT farm tokens.
            </p>

            {error && <AlertView message={error} type="error" />}
            {alert && <AlertView message={alert} type="warning" />}
            {success && <AlertView message={success} type="success" />}
            {isSending && (
              <AlertView
                message="Transaction is being sent to the server"
                type="warning"
              />
            )}
            {chains.length === 0 && (
              <AlertView message="No chains available" type="error" />
            )}

            <form className="space-y-4" onSubmit={handleSubmit}>
              <input
                type="text"
                value={address}
                onChange={(e) => onAddressChanged(e.target.value)}
                placeholder="0xAD..."
                className="w-full px-4 py-[14px] border border-gray-200 rounded-lg font-mono text-base placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500"
              />

              <button
                type="submit"
                disabled={isSending}
                className="w-full bg-[#4477FF] text-white font-medium py-[14px] rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSending ? "Sending..." : "Submit"}
              </button>
            </form>

            <a
              href="https://docs.inco.org/getting-started/connect-metamask"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-[#4477FF] text-sm mt-4 hover:underline font-mono"
            >
              Add Rivest Network to Wallet
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </a>
          </div>
        )}
      </div>

      <div
        className="w-full md:w-1/2 h-64 md:h-full bg-cover bg-center flex items-center justify-center"
        style={{
          backgroundImage: "url('/hero-image.png')",
        }}
      />
    </div>
  );
}

export default FaucetPage;
